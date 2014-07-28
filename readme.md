
# wpgo.go

Command-line tool to interact with WordPress blogs. You can read, publish and administer to some degree.

## Setup

Works with WordPress.com and self-hosted blogs running [Jetpack](https://jetpack.me/). If running Jetpack it requires the JSON API to be enabled, which should be activated by default.

You need to configure with an authorization token and blog id for each of your blogs. To do so, you can authorize using this site, which looks totally sketchy but it is run by a good guy, trust me: https://rest-access.test.apokalyptik.com/

Once you obtain the token and blog id, create ~/.wpgo.conf with the following two parameters:

    token = ABCDEFGH123456
    blog_id = 123456

You can setup multiple blogs, by specifying multiple entries, labeling each with [name] before parameters. The name used in brackets, is what you'll type on the command-line. 

    [mkaz]
    token = ABCDEF123456
    blog_id = 123456

    [tinker]
    token = ZYXWVU987654
    blog_id = 987654


## Commands

**Single Site**

    $ wpgo read
    
    $ wpgo read [post-id]

    $ wpgo post filename.md

    $ wpgo stats

    $ wpgo upload image.jpg

**Multiple sites using blog name**

    $ wpgo mkaz read

    $ wpgo mkaz read [post-id]

    $ wpgo mkaz post filename.md

    $ wpgo mkaz stats

    $ wpgo mkaz upload image.jpg

### Publishing New Article

See [sample-post.md](https://raw.githubusercontent.com/mkaz/wpgo/master/sample-post.md) for format. It is similar in style to Jekyll or other static-site tools by specifying post details in the front matter. You can specify title, date, category, and tags.


## Contributors

You? Pull requests welcome.


## License

[GNU General Public License, version 2](http://www.gnu.org/licenses/gpl-2.0.html)

