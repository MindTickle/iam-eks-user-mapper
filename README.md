# IAM EKS User Mapper

Derived from the works of [ygrene's](https://github.com/ygrene) [IAM EKS Mapper](https://github.com/ygrene/iam-eks-user-mapper)

# What it does
If you are using EKS for your cluster and IAM to manage your AWS users, this tool allows you to create EKS users on the basis of your IAM groups with fine-grained access contro, this tool allows you to create EKS users on the basis of your IAM groups with fine-grained access control.

## Setting up
1. Have an AWS IAM Group with users that you want to have access to your EKS cluster
2. Create a new IAM User with an IAM ReadOnly policy (or) a new IAM role with IAM ReadOnly Policy and capability to assume role.
3.
  - If you are using an IAM Role, add the ARN for it in `kubernetes/deployment.yml` path: `spec.template.metadata.annotations` with annotation `iam.amazonaws.com/role`: `ROLE_ARN`
  - If you are using an IAM User, add environment variables `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` in `kubernetes/deployment.yml` path: `spec.template.spec.containers.0.env`
4. Update the `AWS_REGION` environment variable in `kubernetes/deployment.yml` if you aren't running in `ap-southeast-1` with your EKS cluster
5. Edit the `kubernetes/deployment.yml` `command:` with both the IAM group name you want to provide access to, and the Kubernetes group each user in the group should be mapped to.
6. Finally:
```bash
$ kubectl apply -f kubernetes/
```

## How the command works

Provide comma separated values for iam-k8s mapping with each mapping represented as
<iam>::<k8s-group>.

Example usages
```bash

# To map all your devops IAM Group as system:masters and devs IAM Group as developer
./app --iam-k8s-group=devops::system:masters,devs::developer

# Support for multiple kubernetes roles (use `|` between K8s roles)
# To map all your devops IAM Group as system:masters and devs IAM Group as both developer and manager
./app --iam-k8s-group=devops::system:masters,devs::developer|manager
```

## Planned features

1. Option to give specific IAM Users access as well (will be a union if that user is part of a provided IAM group as well)
2. Improved CLI Experience
