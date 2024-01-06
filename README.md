# devtools

Just for fun.

## Scripts

### Shell Script Git Commit Message Generator

  
  This Shell script, [`generate_commit.sh`](/scripts/generate-commit.sh) , helps automate the
  process of creating commit messages for your Git repository.

  #### Prerequisites

  The script requires that the following programs be installed on your system:
  jq, gum, and git. If any of these packages are not installed, the script
  will return an error and cease operation.

  #### API Key

  The script interfaces with the OpenAI API for some of its functions. As
  such, you need to have `OPENAI_API_KEY` environment variable set with your
  API key from OpenAI. If this variable is not set, the script will return an
  error and cease operation.

  See the [mods documentation](https://github.com/charmbracelet/mods) for more details on that interaction.

  #### Git Staging

  The script operates on files that you have staged for a commit. Run `git add
  <filename>` to stage your changes before running this script.

  #### Usage

  Run the script in your terminal by typing
  `./generate-commit.sh`.

  #### Interaction

  The script will prompt you to summarize your changes, categorize the type of
  changes made along with the area of the repository most affected. Then, it
  composes a commit message with your input.

  You are given the chance to edit the commit message. The message is composed
  of a DESC (description), a SUMMARY of your changes, the TYPE of changes, and
  the SCOPE of the repository that has been affected.

  You can select any part of the generated commit message to edit. Once you
  confirm, the script will commit your changes with the message.

  **NOTE:** This script uses "mods" and "gum" to interactively get user input,
  and OpenAI API is used to process inputs for constructing the commit
  message. It's essential to have internet connectivity while running this
  script as it interacts with the OpenAI API.

  Run the script wisely, as it directly interacts with your Git repository.
  Remember, this is a helper tool and does not replace a diligent and
  thoughtful manual commit process.

  Happy committing!

