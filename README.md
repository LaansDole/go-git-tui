# Useful tools to enhance your productivity
My goal when writing these scripts is to reduce the amount of time 
you have to leave the keyboard. You can say that using `neovim` would be the best
for this purpose. Until you can setup your own `vim` environment, here are my tools
to fill in that gap!

## How to use
1. Clone the repository:
    ```shell
    git clone https://github.com/AnDoLeLongANZ/useful-tools.git
    cd useful-tools
    ```
2. Make sure that you have `make` command
3. From `./useful-tools` run `make execute`
4. In `.zshrc`, aliasing the command that you want to use, for example:
    ```shell
    alias reset='source ~/.zshrc'               # reset: reset .zshrc
    alias edit='vim ~/.zshrc'                   # vimsh: modify .zshrc
    #### Advanced Git Command ####
    alias gcommit="source your/path/to/useful-tools/fuzzy-git/gh-commit.sh"
    alias gadd="source your/path/to/useful-tools/fuzzy-git/gh-add.sh"
    ```
5. Remember to "reset" your `.zshrc` with `source ~/.zshrc` to update these changes!