package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// func main() {
// 	//test
// 	svc := AwsSession()
// 	key, err := AwsGetSSHkey("dfranco", svc)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(key)
//
// }

func AwsSession() *iam.IAM {
	// Create a session to share configuration, and load external configuration.
	sess := session.Must(session.NewSession())

	// Create the service's client with the session.
	return iam.New(sess)
}

func AwsGetSSHkey(user string, session *iam.IAM) (sshkeys []string, err error) {

	var userKeys []string

	sshinput := &iam.ListSSHPublicKeysInput{
		UserName: aws.String(user),
	}

	//Get user sshs
	sshKeyInfo, err := session.ListSSHPublicKeys(sshinput)
	if err != nil {
		return nil, err
	}

	//I need a comment here
	if len(sshKeyInfo.SSHPublicKeys) > 0 {
		for n := 0; n < len(sshKeyInfo.SSHPublicKeys); n++ {
			publicssh := &iam.GetSSHPublicKeyInput{
				Encoding:       aws.String("SSH"),
				SSHPublicKeyId: aws.String(*sshKeyInfo.SSHPublicKeys[n].SSHPublicKeyId),
				UserName:       aws.String(user),
			}

			publickey, err := session.GetSSHPublicKey(publicssh)
			if err != nil {
				return nil, err
			}

			userKeys = append(userKeys, *publickey.SSHPublicKey.SSHPublicKeyBody)
		}
		return userKeys, nil
	}
	return []string{}, nil
}
