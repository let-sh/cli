#!/usr/bin/env bash
# shellcheck disable=SC2059

set -eo pipefail

reset="\033[0m"
red="\033[31m"
green="\033[32m"
yellow="\033[33m"
cyan="\033[36m"
white="\033[37m"

letsh_get_binary() {
  local os
  case "$(uname -s)" in
  Darwin)
    #      printf "${green}> macOS detected$reset\n"
    os='darwin'
    ;;

  Linux)
    #      printf "${green}> Linux detected$reset\n"
    os='linux'
    ;;

#  CYGWIN* | MINGW32* | MSYS* | MINGW*)
#    #      printf "${green}> Windows detected$reset\n"
#    os='win'
#    ;;

  *)
    printf "$red> Unsupported OS.$reset\n"
    exit 1
    ;;
  esac

  printf "$cyan> Downloading binary...$reset\n"
  if [ "$1" = '--version' ]; then
    # Validate that the version matches MAJOR.MINOR.PATCH to avoid garbage-in/garbage-out behavior
    version=${v#$2}
    if echo "$version" | grep -qE "^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$"; then
      url="https://install.let-sh.com/cli_${version}_${os}_amd64.tar.gz"
    else
      printf "$red> Version number invalid.$reset\n"
      exit 1
    fi
  else
    version=$(curl -sS https://install.let-sh.com/version | grep "latest:" | cut -f 2 -d ":")
    url="https://install.let-sh.com/cli_${version}_${os}_amd64.tar.gz"
  fi

#  if [ $os = 'win' ]; then
#    url="$url.exe"
#  fi

  # Get both the tarball and its GPG signature
  # if curl --fail --progress-bar -L -o ~/.let/bin/lets "$url"; then
  tmpfile=$(mktemp /tmp/let-sh-bundle.XXXXXX)
  if curl --fail --progress-bar -L "$url" > $tmpfile; then
    mkdir -p ~/.let/bin
    tar xzvf $tmpfile -C ~/.let/bin/ > /dev/null 2>&1
    printf "$cyan> Chmod-ing file ~/.let/bin/lets...$reset\n"
    chmod +x ~/.let/bin/lets
    rm -f $tmpfile
  else
    rm -f $tmpfile
    printf "$red> Failed to download $url.$reset\n"
    exit 1
  fi
}

letsh_detect_profile() {
  if [ -n "${PROFILE}" ] && [ -f "${PROFILE}" ]; then
    echo "${PROFILE}"
    return
  fi

  local DETECTED_PROFILE
  DETECTED_PROFILE=''
  local SHELLTYPE
  SHELLTYPE="$(basename "/$SHELL")"

  if [ "$SHELLTYPE" = "bash" ]; then
    if [ -f "$HOME/.bashrc" ]; then
      DETECTED_PROFILE="$HOME/.bashrc"
    elif [ -f "$HOME/.bash_profile" ]; then
      DETECTED_PROFILE="$HOME/.bash_profile"
    fi
  elif [ "$SHELLTYPE" = "zsh" ]; then
    DETECTED_PROFILE="$HOME/.zshrc"
  elif [ "$SHELLTYPE" = "fish" ]; then
    DETECTED_PROFILE="$HOME/.config/fish/config.fish"
  fi

  if [ -z "$DETECTED_PROFILE" ]; then
    if [ -f "$HOME/.profile" ]; then
      DETECTED_PROFILE="$HOME/.profile"
    elif [ -f "$HOME/.bashrc" ]; then
      DETECTED_PROFILE="$HOME/.bashrc"
    elif [ -f "$HOME/.bash_profile" ]; then
      DETECTED_PROFILE="$HOME/.bash_profile"
    elif [ -f "$HOME/.zshrc" ]; then
      DETECTED_PROFILE="$HOME/.zshrc"
    elif [ -f "$HOME/.config/fish/config.fish" ]; then
      DETECTED_PROFILE="$HOME/.config/fish/config.fish"
    fi
  fi

  if [ -n "$DETECTED_PROFILE" ]; then
    echo "$DETECTED_PROFILE"
  fi
}

letsh_link() {
  printf "$cyan> Adding to \$PATH...$reset\n"
  LETSH_PROFILE="$(letsh_detect_profile)"
  SOURCE_STR="\nexport PATH=\"\$HOME/.let/bin:\$PATH\"\n"

  if [ -z "${LETSH_PROFILE-}" ]; then
    printf "$red> Profile not found. Tried ${LETSH_PROFILE} (as defined in \$PROFILE), ~/.bashrc, ~/.bash_profile, ~/.zshrc, and ~/.profile.\n"
    echo "> Create one of them and run this script again"
    echo "> Create it (touch ${LETSH_PROFILE}) and run this script again"
    echo "   OR"
    printf "> Append the following lines to the correct file yourself:$reset\n"
    command printf "${SOURCE_STR}"
  else
    if ! grep -q 'let/bin' "$LETSH_PROFILE"; then
      if [[ $LETSH_PROFILE == *"fish"* ]]; then
        # shellcheck disable=SC2016
        command fish -c 'set -U fish_user_paths $fish_user_paths ~/.let/bin'
        printf "$cyan> We've added ~/.let/bin to your fish_user_paths universal variable\n"
      else
        command printf "$SOURCE_STR" >>"$LETSH_PROFILE"
        printf "$cyan> We've added the following to your $LETSH_PROFILE\n"
      fi

      echo "> If this isn't the profile of your current shell then please add the following to your correct profile:"
      printf "   $SOURCE_STR$reset\n"
    fi

    version=$($HOME/.let/bin/lets version) || (
      printf "$red> let.sh was installed, but doesn't seem to be working :(.$reset\n"
      exit 1
    )

    printf "$green> Successfully installed let.sh $version! Reloading shell...$reset\n"
    
    source "$DETECTED_PROFILE"
  fi
}

letsh_reset() {
  unset -f letsh_install letsh_reset letsh_get_binary letsh_link letsh_detect_profile
}

letsh_install() {
  printf "${white}Installing let.sh!$reset\n"

  if [ -d "$HOME/.let/bin" ]; then
    if which lets; then
      local specified_version
      if [ "$1" = '--version' ]; then
        specified_version=$2
      else
        specified_version=$(curl -sS https://install.let-sh.com/version | grep "latest:" | cut -f 2 -d ":")
      fi
      letsh_version=$(lets version)

      if [ "$specified_version" = "$letsh_version" ]; then
        printf "$green> let.sh is already at the $specified_version version.$reset\n"
        exit 0
      else
        printf "$yellow> $letsh_version is already installed, Specified version: $specified_version.$reset\n"
        rm -rf "$HOME/.let/bin"
      fi
    else
      printf "$red> $HOME/.let/bin already exists, possibly from a past let.sh install.$reset\n"
      printf "$red> Remove it (rm -rf $HOME/.let/bin) and run this script again.$reset\n"
      exit 0
    fi
  fi

  letsh_get_binary $1 $2
  letsh_link
  letsh_reset
}

cd ~
letsh_install $1 $2
