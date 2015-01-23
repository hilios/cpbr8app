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
app.config(function($httpProvider) {
  // Serialize the given Object into a key-value pair string. This method
  // expects an object and will default to the toString() method.
  // --
  // NOTE: This is an atered version of the jQuery.param() method which
  // will serialize a data collection for Form posting.
  // --
  // https://github.com/jquery/jquery/blob/master/src/serialize.js#L45
  function _formDataSerializer(data) {
    if (angular.isUndefined(data)) {
      return;
    }
    // If this is not an object, defer to native stringification.
    if (!angular.isObject(data)) {
      return((data == null) ? "" : data.toString());
    }
    var buffer = [];
    // Serialize each key in the object.
    for (var name in data) {
      if (!data.hasOwnProperty(name)) {
        continue;
      }
      var value = data[name];

      buffer.push(encodeURIComponent(name) + "=" +
        encodeURIComponent((value == null) ? "" : value)
      );
    }
    // Serialize the buffer and clean it up for transportation.
    var source = buffer.join("&").replace(/%20/g, "+");

    return source;
  }

  function _transformRequestDataToForm(data, getHeaders) {
    if (data) {
      data = angular.fromJson(data);

      var headers = getHeaders();
      headers["Content-Type"] = "application/x-www-form-urlencoded";

      return _formDataSerializer(data);
    }

    return data;
  }

  $httpProvider.defaults.transformRequest.push(_transformRequestDataToForm);
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

  function getId(model) {
    var id = null;

    if (model.hasOwnProperty('id')) {
      id = model.id;
    }

    return id;
  }

  function all() {
    var url = _urlFor('/tasks');
    return $http.get(url);
  }

  function put(model) {
    var url, q, id = getId(model);


    if (id) {
      url = _urlFor('/task', id);
      q = $http.put(url, model);
    } else {
      url = _urlFor('/task');
      q = $http.post(url, model);
    }

    return q;
  }

  function remove(model) {
    var url, id = getId(model);

    url = _urlFor('/task', id);
    return $http.delete(url)
  }

  return {
    'all': all,
    'put': put,
    'remove': remove,
    'getId': getId
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

  $scope.tasks = null;
  $scope.currentModel = null;

  $scope.open = function(model) {
    $scope.currentModel = model;
  }

  $scope.save = function(model) {
    q = Tasks.put(model)
      .success(function(data) {
        var isNew = $scope.tasks.map(Tasks.getId).indexOf(data.id) == -1;
        if (isNew) {
          $scope.tasks.push(data);
        }
      })
      .finally(function() {
        $scope.close();
      });

    Spinner.watch(q)
  }

  $scope.toggle = function(model) {
    model.ok = !model.ok;
    $scope.save(model);
  }

  $scope.remove = function(model) {
    q = Tasks.remove(model)
      .success(function() {
        // Remove the model from the list
        var i = $scope.tasks.indexOf(model);
        $scope.tasks.splice(i, 1) ;
      });

    Spinner.watch(q)
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
      'onEdit': '&',
      'onRemove': '&',
      'onToggle': '&'
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
