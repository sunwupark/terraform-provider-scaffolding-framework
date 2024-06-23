terraform {
    required_providers {
        sunwupark = {
            source = "registry.terraform.io/study/sunwupark-ossca"
        }
    }   
}

# provider "sunwupark"{
#     username = "sunwupark"
#     password = "1234"
#     host     = "http://localhost:19090"
# }

# resource "sunwupark_friend" "friend" {
#     name = "sunwupark"
#     address = "Seoul"
#     description = "Hello, I'm sunwupark"
#     image = "https://avatars.githubusercontent.com/u/12345678?v=4"
# }