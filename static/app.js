var app = angular.module('tasks', []);
// Configuration
app.value('API_ENDPOINT', 'http://cpbr8app.herokuapp.com')
// Services
app.factory('Tasks', function($http, API_ENDPOINT) {

  function _buildUrl(path, id) {
    url = API_ENDPOINT + path;

    if (id) {
      url = url + '?id=' + id;
    }

    return url;
  }

  function all() {
    var url = _buildUrl('/tasks');
    return $http.get(url);
  }

  return {
    'all': all
  };
});
// Controllers
app.controller('TasksController', function ($scope, Tasks) {
  $scope.tasks = [];

  $scope.refresh = function() {
    Tasks.all().success(function(data) {
      $scope.tasks = data.tasks;
    });
  }

  $scope.refresh();
});
app.controller('TaskController', function ($scope, Tasks) {

});
// Directives
app.directive('task', function() {
  return {
    templateUrl: 'task.html',
    require: 'ngModel',
    scope: {
      'ngModel': '='
    }
  };
});
// Filters
app.filter('len', function() {
  return function(input) {
    try {
      return input.length
    } catch(e) {
      return 0
    }
  }
});
