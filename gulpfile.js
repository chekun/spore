var elixir = require('laravel-elixir');

elixir(function(mix) {
    mix.browserify('app.js');
    mix.copy('resources/assets/logo.png', 'public/assets/logo.png');
    mix.styles(['app.css', 'nprogress.css']);
});
