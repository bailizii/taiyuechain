package cim

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/pkg/errors"
	"math/big"
	"crypto/x509/pkix"
	"time"
	"log"
	"crypto/rand"
	"encoding/base64"
	"os"
	"fmt"
	"github.com/taiyuechain/taiyuechain/crypto/taiCrypto"
	"github.com/taiyuechain/taiyuechain/crypto"
	"crypto/ecdsa"

	"crypto/elliptic"
)

func GetIdentityFromByte(idBytes []byte) (Identity, error) {
	cert, err := GetCertFromPem(idBytes)
	if err != nil {
		return nil, err
	}

	keyImporter := &x509PublicKeyImportOptsKeyImporter{}
	opts := &X509PublicKeyImportOpts{Temporary: true}

	certPubK, err := keyImporter.KeyImport(cert, opts)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed importing key with opts")
	}

	identity, err := NewIdentity(cert, certPubK)
	if err != nil {
		return nil, err
	}
	return identity, nil
}

func GetCertFromPem(idBytes []byte) (*x509.Certificate, error) {
	if idBytes == nil {
		return nil, errors.New("getCertFromPem error: nil idBytes")
	}

	// Decode the pem bytes
	pemCert, _ := pem.Decode(idBytes)
	if pemCert == nil {
		return nil, errors.Errorf("getCertFromPem error: could not decode pem bytes [%v]", idBytes)
	}

	// get a cert
	var cert *x509.Certificate
	cert, err := x509.ParseCertificate(pemCert.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "getCertFromPem error: failed to parse x509 cert")
	}

	return cert, nil
}


func  CreateIdentity(priv string) bool {
	var private taiCrypto.TaiPrivateKey
	//var public taiCrypto.TaiPublicKey
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Country:            []string{"China"},
			Organization:       []string{"Yjwt"},
			OrganizationalUnit: []string{"YjwtU"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	//ecdsa, err := taiCrypto.HexToTaiPrivateKey(priv)
	//var thash taiCrypto.THash
	//caecda, err := private.ToECDSACA(ecdsa.HexBytesPrivate)
	caecda, err := private.ToECDSACA([]byte(priv))
	if err != nil {
		log.Println("create ca failed", err)
		return false
	}
	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, &caecda.Private.PublicKey, &caecda.Private)
	if err != nil {
		log.Println("create ca failed", err)
		return false
	}

	encodeString := base64.StdEncoding.EncodeToString(ca_b)

	fileName := priv[:4] + "ca.pem"
	dstFile, err := os.Create(fileName)
	if err != nil {
		return false
	}
	defer dstFile.Close()
	priv_b, _ := x509.MarshalECPrivateKey(&caecda.Private)
	encodeString1 := base64.StdEncoding.EncodeToString(priv_b)
	if err != nil {
		fmt.Println(err)
	}
	fileName1 := priv[:4] + "ca.key"
	dstFile1, err := os.Create(fileName1)
	if err != nil {
		return false
	}
	defer dstFile1.Close()
	dstFile1.WriteString(encodeString1 + "\n")
	fmt.Println(encodeString)
	return true
}

func  CreateIdentity2(priv , priv2 *ecdsa.PrivateKey,name string) bool {
	//var private taiCrypto.TaiPrivateKey
	//var public taiCrypto.TaiPublicKey
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Country:            []string{"China"},
			Organization:       []string{"Yjwt"},
			OrganizationalUnit: []string{"YjwtU"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	//ecdsa, err := taiCrypto.HexToTaiPrivateKey(priv)
	//var thash taiCrypto.THash
	//caecda, err := private.ToECDSACA(ecdsa.HexBytesPrivate)
	//caecda, err := private.ToECDSACA([]byte(priv))
	//pub := crypto.FromECDSAPub(&priv.PublicKey)
	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, &priv.PublicKey, priv)
	if err != nil {
		log.Println("create ca failed", err)
		return false
	}

	encodeca := base64.StdEncoding.EncodeToString(ca_b)
	fmt.Println(encodeca)
	bytes, _ := base64.StdEncoding.DecodeString(encodeca)
	/*var data []byte
	if strings.Contains(string(bytes), "-BEGIN CERTIFICATE-") {
		block, _ := pem.Decode(ca_b)
		if block == nil {
			fmt.Println("that ca not right")
		}
		data = block.Bytes
	}*/

	//theCert, err := x509.ParseCertificate(data)

	 //t  :="696b0620068602ecdda42ada206f74952d8c305a811599d463b89cfa3ba3bb98"

	theCert, err := x509.ParseCertificate(bytes)
	pubk1 := theCert.PublicKey
	var publicKeyBytes []byte

	switch pub2 := pubk1.(type) {
	case *ecdsa.PublicKey:
		publicKeyBytes = elliptic.Marshal(pub2.Curve, pub2.X, pub2.Y)

	}

	//pkcString := string()
	fmt.Println(publicKeyBytes)
	fmt.Println(crypto.FromECDSAPub(&priv.PublicKey))
	fmt.Println(crypto.FromECDSAPub(&priv2.PublicKey))
	if(string(publicKeyBytes) == string(crypto.FromECDSAPub(&priv2.PublicKey))){
		fmt.Println("1111succes")
	}else{
		if(string(publicKeyBytes) == string(crypto.FromECDSAPub(&priv.PublicKey))){
			fmt.Println("222success")
		}
	}
	//fmt.Println(pkcString)
	/*if(string(crypto.FromECDSA(pkcert)) == (string(crypto.FromECDSA(priv)))){
		fmt.Println("=====")
	}else{
		fmt.Println("not =====")
	}*/


	encodeString := base64.StdEncoding.EncodeToString(ca_b)
	fileName := "./testdata/testcert/"+name + "ca.pem"
	dstFile, err := os.Create(fileName)
	if err != nil {
		return false
	}
	dstFile.WriteString(encodeString +"\n")
	defer dstFile.Close()
	/*
	priv_b, _ := x509.MarshalECPrivateKey(priv)
	encodeString1 := base64.StdEncoding.EncodeToString(priv_b)
	if err != nil {
		fmt.Println(err)
	}
	fileName1 := "test1" + "ca.key"
	dstFile1, err := os.Create(fileName1)
	if err != nil {
		return false
	}
	defer dstFile1.Close()
	dstFile1.WriteString(encodeString1 + "\n")
	fmt.Println(encodeString)*/
	return true
}


