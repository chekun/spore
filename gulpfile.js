var elixir = require('laravel-elixir');

elixir(function(mix) {
    mix.browserify('app.js');
    mix.copy('resources/assets/logo.png', 'public/assets/logo.png');
    mix.copy('node_modules/bootstrap/dist/css/bootstrap.css', 'resources/css/bootstrap.css');
    mix.copy('node_modules/bootstrap/dist/css/bootstrap-theme.css', 'resources/css/bootstrap-theme.css');
    mix.styles(['bootstrap.css', 'bootstrap-theme.css', 'app.css', 'nprogress.css']);
});
