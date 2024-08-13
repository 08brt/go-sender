# Email Sender with Attachment in Go

This application is designed to streamline the process of sending your CV via email when you come across a recruiter's direct email on LinkedIn. Simply input the recruiter's email address, and the application will automatically send a predefined email with your CV attached.

### Save Time with Automated Emails!

## Prerequisites

- **Go** 1.16 or later
- **Gmail Account**: Ensure you have the correct SMTP settings and credentials for your Gmail account.

## Configuration

Before running the application, you need to set the following configurations directly in the `main.go` file:

- **Username**: Your Gmail address (e.g., `username@gmail.com`)
- **Password**: Your Gmail app password (not your regular Gmail password; ensure you have generated an app-specific password if needed).
- **File Path**: Set the file path of your CV in the `main.go` file.

## Usage

1. Clone the repository or download the code to your local machine.
2. Navigate to the directory containing the code.
3. Modify the `main.go` file to include your Gmail credentials and file path.
4. Edit the .bat file with correct path and run the .bat file