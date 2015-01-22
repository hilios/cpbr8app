var app = angular.module('tasks', ['ngAnimate']);
// Constants and variables
app.constant('API_ENDPOINT', 'http://cpbr8app.herokuapp.com')
// Configuration
app.config(function($sceDelegateProvider, API_ENDPOINT) {
  $sceDelegateProvider.resourceUrlWhitelist([
    'self',
    API_ENDPOINT + '/**'
  ]);
});
// Services
app.factory('Tasks', function($http, API_ENDPOINT) {

  function _urlFor(path, id) {
    var url = API_ENDPOINT;

    if (id) {
      url = API_ENDPOINT + path +'?id=' + id;
    } else {
      url = API_ENDPOINT + path;
    }

    return url;
  }

  function all() {
    var url = _urlFor('/tasks');
    return $http.get(url);
  }

  function put(data) {
    var id, url, q;
    if (data.hasOwnProperty('id')) {
      id = data.id;
    }

    if (id) {
      url = _urlFor('/task', id);
      q = $http.put(url, data);
    } else {
      url = _urlFor('/task');
      q = $http.post(url, data);
    }

    return q;
  }

  return {
    'all': all,
    'put': put
  };
});
app.factory('Spinner', function($rootScope, $q) {

  function watch() {
    var promisses = Array.prototype.slice.call(arguments);

    $rootScope.$broadcast('spinner:show');

    $q.all(promisses).finally(function() {
      $rootScope.$broadcast('spinner:hide');
    });
  }

  return {
    'watch': watch
  };
});
// Controllers
app.controller('TasksController', function ($scope, Tasks, Spinner) {
  var q;

  $scope.tasks = [];
  $scope.currentModel = null;

  function _findById(id) {
    return function(model) {
      return model.id == id;
    }
  }

  $scope.open = function(model) {
    $scope.currentModel = model;
  }

  $scope.save = function(model) {
    q = Tasks.put(model)
      .success(function(data) {
        var isNew = !$scope.tasks.some(_findById(data.id));
        if (isNew) {
          $scope.tasks.push(data)
        }
      })
      .finally(function() {
        $scope.close();
      });

    Spinner.watch(q)
  }

  $scope.toggle = function(model) {
    model.ok = !model.ok;
    $scope.save()
  }

  $scope.close = function() {
    $scope.currentModel = null;
  }

  $scope.refresh = function() {
    q = Tasks.all()
      .success(function(data) {
        $scope.tasks = data.tasks;
      });

    Spinner.watch(q);
  }

  $scope.refresh();
});
// Directives
app.directive('task', function() {
  return {
    templateUrl: 'task.html',
    scope: {
      'ngModel': '=',
      'onEdit': '&'
    }
  };
});
app.directive('spinner', function() {
  return {
    templateUrl: 'spinner.html',
    scope: true,
    link: function(scope, el, attrs) {
      scope.isLoading = true;
      scope.$on('spinner:show', function() {
        scope.isLoading = true;
      });
      scope.$on('spinner:hide', function() {
        scope.isLoading = false;
      });
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
