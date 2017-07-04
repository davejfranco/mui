package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func AwsSession() *iam.IAM {
	// Create a session to share configuration, and load external configuration.
	sess := session.Must(session.NewSession())

	// Create the service's client with the session.
	return iam.New(sess)
}

func IamCheckUser(user string, session *iam.IAM) bool {

	input := &iam.GetUserInput{
		UserName: aws.String(user),
	}

	_, err := session.GetUser(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return false
			case iam.ErrCodeServiceFailureException:
				fmt.Println(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	return true
}

func AwsGetSSHkey(user string, session *iam.IAM) (sshkeys []string, err error) {

	var userKeys []string

	sshinput := &iam.ListSSHPublicKeysInput{
		UserName: aws.String(user),
	}

	//Get user sshs
	sshKeyInfo, err := session.ListSSHPublicKeys(sshinput)
	if err != nil {
		fmt.Println(err)
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

func iamUser(user string) (usr User, err error) {
	//create user based on IAM
	//aws session
	svc := AwsSession()

	userkeys, err := AwsGetSSHkey(user, svc)
	if err != nil {
		return User{}, err
	}

	u := User{
		Name:      user,
		Publickey: userkeys,
	}
	return u, nil
}
