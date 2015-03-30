# gopass
[Pass](http://www.passwordstore.org/) implementation in Go.

Password management should be simple and follow Unix philosophy. With ``gopass``, each password lives inside of a ``gpg`` encrypted file whose filename is the title of the website or resource that requires the password. These encrypted files may be organized into meaningful folder hierarchies, copied from computer to computer, and, in general, manipulated using standard command line file management utilities.

``gopass`` makes managing these individual password files extremely easy. All passwords live in ``~/.password-store``, and gopass provides some nice commands for adding, editing, generating, and retrieving passwords. It's capable of temporarily putting passwords on your clipboard and tracking password changes using git.
