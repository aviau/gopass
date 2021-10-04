# gopass

[![GoDoc](https://godoc.org/github.com/aviau/gopass?status.svg)](http://godoc.org/github.com/aviau/gopass)


[Pass](http://www.passwordstore.org/) implementation in Go.

Password management should be simple and follow Unix philosophy. With ``gopass``, each password lives inside of a ``gpg`` encrypted file whose filename is the title of the website or resource that requires the password. These encrypted files may be organized into meaningful folder hierarchies, copied from computer to computer, and, in general, manipulated using standard command line file management utilities.

``gopass`` makes managing these individual password files extremely easy. All passwords live in ``~/.password-store``, and gopass provides some nice commands for adding, editing, generating, and retrieving passwords. It's capable of temporarily putting passwords on your clipboard and tracking password changes using git.

## Install

gopass is available in official Debian repositories. Install it with ``apt-get install gopass``.

## Project Status

This section was just added so that I could get an idea of where I am at.

### ``gopass init``

- [X] Creates a folder and a .gpg-id file
- [X] Support ``--path`` option
- [X] Support multiple GPG ids
- [X] Re-encryption functionality
- [ ] Should output: ``Password store initialized for [gpg-id].``
- [ ] ``--clone <url>`` allows to init from an existing repo

### ``gopass insert``

- [X] ``gopass insert test.com`` prompts for a password and creates a test.com.gpg file
- [X] Multi-line support
- [X] Create a git commit
- [X] Prompt before overwriting an existing password, unless --force or -f is specified.
- [ ] When inserting in a folder with a .gpg-id file, insert should use the .gpg-id file's key

### ``gopass show``

- [X] ``gopass show test.com`` will display the content of test.com.gpg
- [X] ``--clip, -c`` copies the first line to the clipboard
- [ ] ``--clip, -c`` clears after a while
- [ ] ``--password``, and ``--username`` options.

Accepted format:
```
<the_password>
login: <the_login>
url: <the_url>
```

### ``gopass connect`` (or ``ssh``)

This new command should connect to a server using an encrypted rsa key.

### ``gopass ls``

- [X] ``gopass ls`` shows the content of the password store with ``tree``
- [X] ``gopass`` invokes ``gopass ls`` by default
- [X] ``gopass ls subfolder`` calls tree on the subfolder only
- [ ] Hide .gpg at the end of each entry
- [X] First output line should be ``Password Store``

### ``gopass rm``

- [X] ``gopass rm test.com`` removes the test.com.gpg file
- [X] ``gopass remove`` and ``gopass delete`` aliases
- [X] ``gopass rm -r folder`` (or ``--recursive``)  will remove a folder and all of it's content (not interactive!)
- [X] Ask for confirmation

### ``gopass find``

- [X] ``gopass find python.org test`` will show a tree with password entries that match python.org or test
- [X] Accepts one or many search terms

### ``gopass cp``

- [X] ``gopass cp old-path new-pah`` copies a password to a new path
- [X] Dont overwrite

### ``gopass mv``

- [X] ``gopass mv old-path new-path`` moves a password to a new path
- [X] Dont overwrite

### ``gopass git``

- [X] Pass commands to git
- [ ] ``gopass git init`` should behave differently with an existing password store
- [ ] Add tests

### ``gopass edit``

- [X] ``gopass edit test.com`` will open a text editor and let you edit the password

### ``gopass grep``

- [X] ``gopass grep searchstring`` will search for the given string inside all of the encrypted passwords


### ``gopass generate``

- [X] ``gopass generate [pass-name] [pass-length]`` Genrates a new password using of length pass-length and inserts it into pass-name.
- [X] ``--no-symbols, -n``
- [ ] ``--clip, -c``
- [ ] ``--in-place, -i``
- [X] ``--force, -f``
- [X] Prompt before overwriting an existing password, unless --force or -f is specified.

## Note

- This isn't [gopass.pw](https://www.gopass.pw/).
