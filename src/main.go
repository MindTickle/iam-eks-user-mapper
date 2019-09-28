package main

import (
	"flag"
	"github.com/kataras/golog"
	"gopkg.in/yaml.v2"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// generateUserRoles returns a map of usernames and their UserRoles
func generateUserRoles(iamK8sGroups []string) map[string]UserRoles {
	userRoles := make(map[string]UserRoles)

	// For each iam, extract users and map them to their k8s roles
	for _, iamK8sGroup := range iamK8sGroups {
		iam, userK8sRoles := extractIAMK8sFromString(iamK8sGroup)
		users := getAwsIamGroup(iam)
		for _, user := range users.Users {
			if _, exists := userRoles[*user.UserName]; !exists {
				userRoles[*user.UserName] = UserRoles{IAMArn: *user.Arn, IAMUsername: *user.UserName, K8sRoles: []string{}}
			}
			userRoles[*user.UserName] = userRoles[*user.UserName].SetK8sRoles(strings.Split(userK8sRoles, "|"))
		}
	}
	// Remove duplicated roles
	for iamUsername := range userRoles {
		userRoles[iamUsername] = userRoles[iamUsername].UniqueK8sRoles()
	}
	return userRoles
}

func main() {
	flagUsage := `
		Provide comma separated values for iam-k8s mapping with each mapping represented as
		<iam>::<k8s-group>.
		Example usage
		--iam-k8s-group=devops::system:masters,devs::developer|manager
	`
	iamK8sGroupRaw := flag.String("iam-k8s-group", "", flagUsage)
	flag.Parse()

	iamK8sGroups := strings.Split(*iamK8sGroupRaw, ",")

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		userRoles := generateUserRoles(iamK8sGroups)
		cf, err := clientset.CoreV1().ConfigMaps("kube-system").Get("aws-auth", metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}
		var newConfig []MapUserConfig

		for _, userRole := range userRoles {
			newConfig = append(newConfig, MapUserConfig{
				UserArn:  userRole.IAMArn,
				Username: userRole.IAMUsername,
				Groups:   userRole.K8sRoles,
			})
		}

		// If there are no users to add, the config map will be empty
		// Since this will never be the intended purpose of the user
		// and this case is more likely to happen due to a bug
		// we'll just skip the changes
		if len(newConfig) == 0 {
			golog.Info("No users found, config will not be changed")
			continue
		}

		roleStr, err := yaml.Marshal(newConfig)
		if err != nil {
			golog.Error(err)
		}
		cf.Data["mapUsers"] = string(roleStr)

		newCF, err := clientset.CoreV1().ConfigMaps("kube-system").Update(cf)
		if err != nil {
			golog.Error(err)
		} else {
			golog.Info("Successfully updated user roles")
			golog.Info(newCF)
		}
		time.Sleep(1 * time.Minute)
	}
}
