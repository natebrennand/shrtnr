
/*
 *  Navbar Controller
 */

angular.module('app.navbar', [])
.controller('NavbarController', function ($scope, $location) {

  $scope.currentTab = $location.path();

  $scope.$on("$routeChangeStart", function (event, next, current) {
    $scope.currentTab = $location.path();
  });

});

