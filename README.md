# gopass
[Pass](http://www.passwordstore.org/) implementation in Go.

Password management should be simple and follow Unix philosophy. With ``gopass``, each password lives inside of a ``gpg`` encrypted file whose filename is the title of the website or resource that requires the password. These encrypted files may be organized into meaningful folder hierarchies, copied from computer to computer, and, in general, manipulated using standard command line file management utilities.

``gopass`` makes managing these individual password files extremely easy. All passwords live in ``~/.password-store``, and gopass provides some nice commands for adding, editing, generating, and retrieving passwords. It's capable of temporarily putting passwords on your clipboard and tracking password changes using git.

## Project Status

This section was just added so that I could get an idea of where I am at. I won't check any feature unless there is a corresponding test. For now, **it is therefore an understatement of what is done and what isn't**.

###``gopass init``

- [ ] Creates a folder and a .gpg-id file
- [ ] Support ``--path`` option
- [ ] Re-encryption functionality
- [ ] Should output: ``Password store initialized for [gpg-id].``
- [ ] ``--clone <url>`` allows to init from an existing repo

###``gopass insert``

- [ ] ``gopass insert test.com`` prompts for a password and creates a test.com.gpg file
- [ ] Multi-line support
- [ ] Create a git commit
- [ ] When inserting in a folder with a .gpg-id file, insert should use the .gpg-id file's key

###``gopass show``

- [ ] ``gopass show test.com`` will display the content of test.com.gpg
- [ ] ``--clip, -c`` copies the first line to the clipboard
- [ ] ``--password``, and ``--username`` options.

Accepted format:
```
<the_password>
login: <the_login>
url: <the_url> 
```

###``gopass connect`` (or ``ssh``)

This new command should connect to a server using an encrypted rsa key. 

###``gopass ls``

- [ ] ``gopass ls`` shows the content of the password store with ``tree``
- [ ] ``gopass`` invokes ``gopass ls`` by default
- [ ] ``gopass ls subfolder`` calls tree on the subfolder only
- [ ] Hide .gpg at the end of each entry
- [ ] Accept subfolder argument
- [ ] First output line should be ``Password Store``

###``gopass rm``

- [ ] ``gopass rm test.com`` removes the test.com.gpg file
- [ ] ``gopass remove`` and ``gopass delete`` aliases
- [ ] ``gopass rm -r folder`` (or ``--recursive``)  will remove a folder and all of it's content (not interactive!)
- [ ] Ask for confirmation

###``gopass find``

- [ ] ``gopass find python.org test`` will show a tree with password entries that match python.org or test
- [ ] Accepts one or many search terms

###``gopass cp``

- [ ] ``gopass cp old-path new-pah`` copies a password to a new path
- [ ] Dont overwrite

###``gopass mv``

- [ ] ``gopass mv old-path new-path`` moves a password to a new path
- [ ] Dont overwrite

###``gopass git``

- [ ] Pass commands to git
- [ ] ``gopass git init`` should behave differently with an existing password store
- [ ] Add tests

###``gopass edit``

- [ ] ``gopass edit test.com`` will open a text editor and let you edit the password

###``gopass grep``

- [ ] ``gopass grep searchstring`` will search for the given string inside all of the encrypted passwords


``gopass generate``
-------------------
- [ ] ``gopass generate [pass-name] [pass-length]`` Genrates a new password using of length pass-length and inserts it into pass-name.
- [ ] ``--no-symbols, -n``
- [ ] ``--clip, -c``
- [ ] ``--in-place, -i``
- [ ] ``--force, -f``
