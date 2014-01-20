[godoc.org]: http://godoc.org/github.com/bmatsuo/gonew "godoc.org"

#Migrating from Gonew Classic to v2

##Config migration

The configuration has changed radically. You will definitely want to start with
the `gonew.json.example` file provided and overwrite the relevant parts with
information from your `~/.gonewrc` file. So, first, copy the example file to the
new location Gonew looks for configuration.

    mkdir -p gonew
    cp gonew.json.example ~/.config/gonew.json

Now open both `~/gonew.json` and `~/.config/gonew.json` in a text editor (such as
Sublime Text)

    /Applications/Sublime\ Text.app/Contents/Sublime\ Text ~/.gonewrc ~/.config/gonew.json

You will want to copy your name and email from the rc file to the json config at
`$.Environments.default.User.*`. You can delete either of the json values if you
want.

The default config provides an environment that extends the default, "work". Feel
free to remove this if you don't have need for another root import path to
publish packages from.

##Import path migration

The `host` and `hostuser` settings of the rc file are not directly represented
in json configuration. Instead, the base import path is specified in the
environment configuration. If automated github or google code integration is
something you want, there are adequate tools to build 'hooks' for those tasks.

##Repo migration

Repository management is outside of Genew's concerns in v2. The project template
system provides all the necessary machinery to provide completely custom
solutions.

The default configuration provides git integration as was provided in Classic.
You may want to remove this. If so, remove `$.Projects.git` and all references
to "git" in `$.Projects.*.Inherits`.


##License migration

Licenses are outside of Gonew's concern in v2. The template system is general
enough to completely handle that task. The example configuration shows how
the newbsd license can be instrumented with the project template system.

There are no other license templates in Gonew v2 (yet). If you wanted other
license templates you would have to use `$.ExternalTemplates` in the json config
file. For more details, see `godoc.org`.

##Custom template migration

If you have custom go templates defined for Classic you will need to migrate
those over yourself. Look at the `*.t2` templates in this repository as a guide.
