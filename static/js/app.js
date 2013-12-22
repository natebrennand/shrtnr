
var app = angular.module('app', [
  'app.home',
  'app.navbar',
  'ngRoute'
]);


app.config(['$routeProvider', function($routeProvider) {
  // declaring routes
  $routeProvider
    .when('/about', {
      templateUrl: 'static/partials/about.html'
    })
    .when('/contact', {
      templateUrl: 'static/partials/contact.html'
    })
    .when('/', {
      templateUrl: 'static/partials/home.html',
      controller: 'HomeController'
    })
    .otherwise({
      redirectTo: '/'
    });

}]);

