terraform {
    required_providers {
        aws = {
            source  = "hashicorp/aws"
            version = "~> 4.47.0"
        }
    }

    # Note that in production configurations the
    # terraform state should be stored and loaded from 
    # and S3 bucket configuration.
    #
    # For the purposes of this local demo the terraform state
    # is held on the local file system.
    backend "s3" {
        bucket         = "nathanbrophy-tf-state-demo"
        key            = "state.tfplan"
        region         = "us-west-1"
    }

    required_version = "~> 1.3"
}
