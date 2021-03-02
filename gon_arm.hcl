# The path follows a pattern
# ./dist/BUILD-ID_TARGET/BINARY-NAME
source = ["./dist/lets_macos_arm_darwin_arm64/lets"]
bundle_id = "com.oasis-networks.cli"

apple_id {
  username = "liangzhib@163.com"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Oasis Networks, Inc."
}

dmg{
  output_path= "./dist/lets_macos_arm_darwin_arm64/lets.dmg"
  volume_name= "let.sh cli"
}

zip {
  output_path = "./dist/lets_macos_arm_darwin_arm64/lets.zip"
}