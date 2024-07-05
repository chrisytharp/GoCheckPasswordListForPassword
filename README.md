# GoCheckPasswordListForPassword
checks to see if a password is in a list 

--- Steps to Run the Program ---
1. Create a Directory for Your Project
  mkdir -p ~/go/src/checkpword
  cd ~/go/src/checkpword

2. Initialize a New Go Module
  go mod init checkpword
  Create the checkpword.go File

3. Install the Required Dependencies
  go get github.com/bits-and-blooms/bloom/v3
  go get github.com/schollz/progressbar/v3

4. Run the Program
<h>
go run checkpword.go "your_password_here" "passwordlist.txt"
