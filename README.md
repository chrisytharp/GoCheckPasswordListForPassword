# GoCheckPasswordListForPassword

## Description

This Go program checks if a password is in a list.

## Steps to Run the Program

1. Create a Directory for Your Project
    ```sh
    mkdir -p ~/go/src/checkpword
    cd ~/go/src/checkpword
    ```

2. Initialize a New Go Module
    ```sh
    go mod init checkpword
    ```

3. Create the `checkpword.go` File

4. Install the Required Dependencies
    ```sh
    go get github.com/bits-and-blooms/bloom/v3
    go get github.com/schollz/progressbar/v3
    ```

5. Run the Program
    ```sh
    go run checkpword.go "your_password_here" "passwordlist.txt"
    ```

## Checking if a Password List Contains a Substring

Use this script and command to check if the password list contains a specified substring and print out what passwords contain that substring:
```sh
go run DoesItCoontainSubString.go "substringToCheck" "passwordlist.txt"

