# GoCheckPasswordListForPassword

<h1> checks to see if a password is in a list </h1>

--- Steps to Run the Program ---
1. Create a Directory for Your Project
<h>
mkdir -p ~/go/src/checkpword
<h>
cd ~/go/src/checkpword

3. Initialize a New Go Module
<h>
go mod init checkpword
<h>
Create the checkpword.go File

5. Install the Required Dependencies
<h>
go get github.com/bits-and-blooms/bloom/v3
<h>
go get github.com/schollz/progressbar/v3

7. Run the Program
<h>
go run checkpword.go "your_password_here" "passwordlist.txt"

<h2> Use this to check what if password list contains substring </h2>
<h></h>
go run DoesItCoontainSubString.go "substringToCheck" "passwordlist.txt"
