package provider

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/fetch"
	"github.com/akmalfairuz/finance/module/hash"
	"net/http"
	"time"
)

type AWSCredentials struct {
	AccessKey string `env:"AWS_ACCESS_KEY"`
	SecretKey string `env:"AWS_SECRET_KEY"`
	Region    string `env:"AWS_REGION"`
}

func (cred AWSCredentials) GetHeaders(method, path, service, query string) fetch.Map {
	host := service + "." + cred.Region + ".amazonaws.com"
	now := time.Now().UTC()
	datetime := now.Format("20060102T150405Z")
	amzDate := now.Format("20060102")

	const awsAlgo = "AWS4-HMAC-SHA256"

	canonicalHeaders := "host:" + host + "\n" + "x-amz-date:" + datetime + "\n"
	const signedHeaders = "host;x-amz-date"
	payloadHash := hash.Sha256("")

	canonicalRequest := method + "\n" + path + "\n" + query + "\n" + canonicalHeaders + "\n" + signedHeaders + "\n" + payloadHash

	credentialScope := amzDate + "/" + cred.Region + "/" + service + "/aws4_request"

	strToSign := awsAlgo + "\n" + datetime + "\n" + credentialScope + "\n" + hash.Sha256(canonicalRequest)

	kDate := hash.HmacSha256Bytes([]byte(amzDate), []byte("AWS4"+cred.SecretKey))
	kRegion := hash.HmacSha256Bytes([]byte(cred.Region), kDate)
	kService := hash.HmacSha256Bytes([]byte(service), kRegion)
	signingKey := hash.HmacSha256Bytes([]byte("aws4_request"), kService)

	signature := hash.HmacSha256(strToSign, string(signingKey))

	authHeader := awsAlgo + " Credential=" + cred.AccessKey + "/" + credentialScope + ", SignedHeaders=" + signedHeaders + ", Signature=" + signature

	headers := fetch.Map{}
	headers.Set("Host", host)
	headers.Set("X-Amz-Date", datetime)
	headers.Set("Authorization", authHeader)

	return headers
}

type AmazonSESEmailProvider struct {
	Credentials AWSCredentials
}

var _ EmailProvider = &AmazonSESEmailProvider{}

func NewAmazonSESEMailProvider(credentials AWSCredentials) *AmazonSESEmailProvider {
	return &AmazonSESEmailProvider{
		Credentials: credentials,
	}
}

func (p *AmazonSESEmailProvider) SendEmail(to []string, subject, body, from string) error {
	query := fetch.Map{}
	query.Set("Action", "SendEmail")
	for i, t := range to {
		query.Set(fmt.Sprintf("Destination.ToAddresses.member.%d", i+1), t)
	}
	query.Set("Source", from)
	query.Set("Message.Subject.Data", subject)
	query.Set("Message.Body.Html.Data", body)
	param := query.ToQueryStringAWSStandard()

	headers := p.Credentials.GetHeaders("GET", "/", "email", param)

	resp, err := fetch.Get("https://email."+p.Credentials.Region+".amazonaws.com/?"+param, &fetch.Options{
		Headers: headers,
		Body:    []byte{},
	})
	if err != nil {
		return err
	}
	if resp.Status != http.StatusOK {
		return fmt.Errorf("error status code: %d, response: %s", resp.Status, string(resp.Body))
	}
	return nil
}
