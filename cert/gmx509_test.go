package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/taiyuechain/taiyuechain/cert/crypto/sm2"
)

func TestX500Name(t *testing.T) {
	name := new(pkix.Name)
	name.CommonName = "ID=Mock Root CA"
	name.Country = []string{"CN"}
	name.Province = []string{"Beijing"}
	name.Locality = []string{"Beijing"}
	name.Organization = []string{"org.zz"}
	name.OrganizationalUnit = []string{"org.zz"}
	fmt.Println(name.String())
}

func TestCreateCertificateRequest(t *testing.T) {
	pri, pub, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	sanContents, err := marshalSANs([]string{"foo.example.com"}, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   "test.example.com",
			Organization: []string{"Σ Acme Co"},
		},
		DNSNames: []string{"test.example.com"},

		// An explicit extension should override the DNSNames from the
		// template.
		ExtraExtensions: []pkix.Extension{
			{
				Id:    oidExtensionSubjectAltName,
				Value: sanContents,
			},
		},
	}

	derBytes, err := CreateCertificateRequest(&template, pub, pri, nil)
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("sample.csr", derBytes, 0644)

	csr, err := ParseCertificateRequest(derBytes)
	if err != nil {
		t.Fatal(err)
	}
	csrPub := csr.PublicKey.(*sm2.PublicKey)
	if !bytes.Equal(pub.GetUnCompressBytes(), csrPub.GetUnCompressBytes()) {
		t.Fatal("public key not equals")
	}

	b, err := VerifyDERCSRSign(derBytes, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Fatal("Verify CSR sign not pass")
	}
}

func TestCreateCertificate(t *testing.T) {
	pri, pub, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	sanContents, err := marshalSANs([]string{"foo.example.com"}, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   "test.example.com",
			Organization: []string{"Σ Acme Co"},
		},
		DNSNames: []string{"test.example.com"},

		// An explicit extension should override the DNSNames from the
		// template.
		ExtraExtensions: []pkix.Extension{
			{
				Id:    oidExtensionSubjectAltName,
				Value: sanContents,
			},
		},
	}

	derBytes, err := CreateCertificateRequest(&template, pub, pri, nil)
	if err != nil {
		t.Fatal(err)
	}

	csr, err := ParseCertificateRequest(derBytes)
	if err != nil {
		t.Fatal(err)
	}

	testExtKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	testUnknownExtKeyUsage := []asn1.ObjectIdentifier{[]int{1, 2, 3}, []int{2, 59, 1}}
	cerTemplate := x509.Certificate{
		// SerialNumber is negative to ensure that negative
		// values are parsed. This is due to the prevalence of
		// buggy code that produces certificates with negative
		// serial numbers.
		SerialNumber: big.NewInt(-1),
		NotBefore:    time.Now(),
		NotAfter:     time.Unix(time.Now().Unix()+100000000, 0),

		SubjectKeyId: []byte{1, 2, 3, 4},
		KeyUsage:     x509.KeyUsageCertSign,

		ExtKeyUsage:        testExtKeyUsage,
		UnknownExtKeyUsage: testUnknownExtKeyUsage,

		BasicConstraintsValid: true,
		IsCA:                  true,

		OCSPServer:            []string{"http://ocsp.example.com"},
		IssuingCertificateURL: []string{"http://crt.example.com/ca1.crt"},

		PolicyIdentifiers: []asn1.ObjectIdentifier{[]int{1, 2, 3}},

		CRLDistributionPoints: []string{"http://crl1.example.com/ca1.crl", "http://crl2.example.com/ca1.crl"},
	}

	FillCertificateTemplateByCSR(&cerTemplate, csr)

	cinfo, err := CreateCertificateInfo(&cerTemplate, &cerTemplate, csr)
	if err != nil {
		t.Fatal(err)
	}

	sign, err := sm2.Sign(pri, nil, cinfo.Raw)
	if err != nil {
		t.Fatal(err)
	}

	cer, err := CreateCertificate(cinfo, sign)
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("sample.cer", cer, 0644)

	certificate, err := ParseCertificate(cer)
	if err != nil {
		t.Fatal(err)
	}

	res := CheckSignature(certificate)
	if res != nil{
		t.Fatal("CheckSignature not pass")

	}else{
		fmt.Println("true")
	}

	fmt.Println(certificate.DNSNames)
}

func TestCheckSignatureFrom(t *testing.T) {
	pri, pub, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	sanContents, err := marshalSANs([]string{"foo.example.com"}, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   "test.example.com",
			Organization: []string{"Σ Acme Co"},
		},
		DNSNames: []string{"test.example.com"},

		// An explicit extension should override the DNSNames from the
		// template.
		ExtraExtensions: []pkix.Extension{
			{
				Id:    oidExtensionSubjectAltName,
				Value: sanContents,
			},
		},
	}

	derBytes, err := CreateCertificateRequest(&template, pub, pri, nil)
	if err != nil {
		t.Fatal(err)
	}

	csr, err := ParseCertificateRequest(derBytes)
	if err != nil {
		t.Fatal(err)
	}

	testExtKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	testUnknownExtKeyUsage := []asn1.ObjectIdentifier{[]int{1, 2, 3}, []int{2, 59, 1}}
	cerTemplate := x509.Certificate{
		// SerialNumber is negative to ensure that negative
		// values are parsed. This is due to the prevalence of
		// buggy code that produces certificates with negative
		// serial numbers.
		SerialNumber: big.NewInt(-1),
		NotBefore:    time.Now(),
		NotAfter:     time.Unix(time.Now().Unix()+100000000, 0),

		SubjectKeyId: []byte{1, 2, 3, 4},
		KeyUsage:     x509.KeyUsageCertSign,

		ExtKeyUsage:        testExtKeyUsage,
		UnknownExtKeyUsage: testUnknownExtKeyUsage,

		BasicConstraintsValid: true,
		IsCA:                  true,

		OCSPServer:            []string{"http://ocsp.example.com"},
		IssuingCertificateURL: []string{"http://crt.example.com/ca1.crt"},

		PolicyIdentifiers: []asn1.ObjectIdentifier{[]int{1, 2, 3}},

		CRLDistributionPoints: []string{"http://crl1.example.com/ca1.crl", "http://crl2.example.com/ca1.crl"},
	}

	FillCertificateTemplateByCSR(&cerTemplate, csr)

	cinfo, err := CreateCertificateInfo(&cerTemplate, &cerTemplate, csr)
	if err != nil {
		t.Fatal(err)
	}

	sign, err := sm2.Sign(pri, nil, cinfo.Raw)
	if err != nil {
		t.Fatal(err)
	}

	cer, err := CreateCertificate(cinfo, sign)
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("sample.cer", cer, 0644)

	parentcertificate, err := ParseCertificate(cer)
	if err != nil {
		t.Fatal(err)
	}

	////////////////////////////////////////son
	_, pub_son, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	sanContents_son, err := marshalSANs([]string{"foo.example.com"}, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	template_son := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   "test.example.com",
			Organization: []string{"Σ Acme Co"},
		},
		DNSNames: []string{"test.example.com"},

		// An explicit extension should override the DNSNames from the
		// template.
		ExtraExtensions: []pkix.Extension{
			{
				Id:    oidExtensionSubjectAltName,
				Value: sanContents_son,
			},
		},
	}

	derBytes_son, err := CreateCertificateRequest(&template_son, pub_son, pri, nil)
	if err != nil {
		t.Fatal(err)
	}

	csr_son, err := ParseCertificateRequest(derBytes_son)
	if err != nil {
		t.Fatal(err)
	}

	//testExtKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	//testUnknownExtKeyUsage := []asn1.ObjectIdentifier{[]int{1, 2, 3}, []int{2, 59, 1}}
	cerTemplate_son := x509.Certificate{
		// SerialNumber is negative to ensure that negative
		// values are parsed. This is due to the prevalence of
		// buggy code that produces certificates with negative
		// serial numbers.
		SerialNumber: big.NewInt(-1),
		NotBefore:    time.Now(),
		NotAfter:     time.Unix(time.Now().Unix()+100000000, 0),

		SubjectKeyId: []byte{1, 2, 3, 4},
		KeyUsage:     x509.KeyUsageCertSign,

		ExtKeyUsage:        testExtKeyUsage,
		UnknownExtKeyUsage: testUnknownExtKeyUsage,

		BasicConstraintsValid: true,
		IsCA:                  true,

		OCSPServer:            []string{"http://ocsp.example.com"},
		IssuingCertificateURL: []string{"http://crt.example.com/ca1.crt"},

		PolicyIdentifiers: []asn1.ObjectIdentifier{[]int{1, 2, 3}},

		CRLDistributionPoints: []string{"http://crl1.example.com/ca1.crl", "http://crl2.example.com/ca1.crl"},
	}

	FillCertificateTemplateByCSR(&cerTemplate_son, csr_son)

	cinfo_son, err := CreateCertificateInfo(&cerTemplate_son, parentcertificate, csr_son)
	if err != nil {
		t.Fatal(err)
	}

	sign_son, err := sm2.Sign(pri, nil, cinfo_son.Raw)
	if err != nil {
		t.Fatal(err)
	}

	cer_son, err := CreateCertificate(cinfo_son, sign_son)
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("sample.cer", cer, 0644)

	parentcertificate_son, err := ParseCertificate(cer_son)
	if err != nil {
		t.Fatal(err)
	}


	err = CheckSignatureFrom(parentcertificate_son,parentcertificate);
	if err != nil{
		t.Fatal(err)
	}
}



func TestOpenCertificate(t *testing.T) {
	/*rootPath := "../root.pem"
	rootByte, _ := ReadPemFileByPath(rootPath)*/

	/*certificate, err := ParseCertificate(rootByte)
	if err!=nil{
		t.Fatalf("111")
	}*/


	caPath := "CA.pem"
	cabyte, _ := ReadPemFileByPath(caPath)

	ca, err := ParseCertificate(cabyte)
	if err!=nil{
		t.Fatalf("111")
	}

	/*rootPath := "../root.pem"
	rootbyte, _ := ReadPemFileByPath(rootPath)

	root, err := ParseCertificate(rootbyte)
	if err!=nil{
		t.Fatalf("331")
	}*/

	fmt.Println(ca)
	//st := "n1"
	//by := []byte(st)
	//_,err:=VerifyDERCSRSign(rootByte,[]byte{'1','2','3','4','5','6'})
	//pub := root.PublicKey.(*sm2.PublicKey)
	pub := ca.PublicKey.(*sm2.PublicKey)

	res :=sm2.Verify(pub, nil, ca.RawTBSCertificate, ca.Signature)
	if !res{
		t.Fatalf("222")
	}



}

func TestOpenParentCertificate(t *testing.T) {
	/*rootPath := "../root.pem"
	rootByte, _ := ReadPemFileByPath(rootPath)*/

	/*certificate, err := ParseCertificate(rootByte)
	if err!=nil{
		t.Fatalf("111")
	}*/


	caPath := "CA.pem"
	cabyte, _ := ReadPemFileByPath(caPath)

	ca, err := ParseCertificate(cabyte)
	if err!=nil{
		t.Fatalf("111")
	}

	sitePath := "site.pem"
	sitebyte, _ := ReadPemFileByPath(sitePath)

	site, err := ParseCertificate(sitebyte)
	if err!=nil{
		t.Fatalf("331")
	}

	fmt.Println(ca)
	//st := "n1"
	//by := []byte(st)
	//_,err:=VerifyDERCSRSign(rootByte,[]byte{'1','2','3','4','5','6'})
	//pub := root.PublicKey.(*sm2.PublicKey)
	pub := ca.PublicKey.(*sm2.PublicKey)

	res :=sm2.Verify(pub, nil, site.RawTBSCertificate, site.Signature)
	if !res{
		t.Fatalf("222")
	}



}
/*
func ReadPemFileByPath(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf(
			"Unable to read test certificate from %q - %q "+
				"Does a unit test have an incorrect test file name?\n",
			path, err))
	}

	if strings.Contains(string(data), "-BEGIN CERTIFICATE-") {
		block, _ := pem.Decode(data)
		if block == nil {
			panic(fmt.Sprintf(
				"Failed to PEM decode test certificate from %q - "+
					"Does a unit test have a buggy test cert file?\n",
				path))
		}
		data = block.Bytes
	}
	return data, nil
}*/

