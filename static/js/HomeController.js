/*
 *  Home Controller
 */

angular.module('app.home', [])
.controller('HomeController', function ($scope, $http, $location) {

  $scope.success = false;
  $scope.warning = false;
  $scope.error = false;

  /*  Flashes the message
   *  @param status must be ['success', 'warning', 'error']
   */
  var flash = function (message, status) {
    $scope.success = false;
    $scope.warning = false;
    $scope.error = false;
    $scope.message = "";
    
    $scope[status] = true
    $scope.message = message;
  };

  $scope.shortenURL = function () {
    var data = {
      "LongURL": $scope.longURL,
      "RequestedURL": $scope.shortURL
    };
    $http.post('/', data)
      .success(function (resp) {
        flash("Your shortened URL is " +
          $location.absUrl().replace("#/", "") + resp.URL,
          "success");
        console.log(resp);
      })
      .error(function (err) {
        flash(err, 'error');
        console.log(err);
      });
  };

});

