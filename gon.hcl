# The path follows a pattern
# ./dist/BUILD-ID_TARGET/BINARY-NAME
source = ["./dist/lets_macos_darwin_amd64/let"]
bundle_id = "com.oasis-networks.cli"

apple_id {
  username = "@env:AC_USERNAME"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Oasis Networks, Inc."
}

dmg{
  output_path= "./dist/lets_macos_darwin_amd64/let.dmg"
  volume_name= "let.sh cli"
}