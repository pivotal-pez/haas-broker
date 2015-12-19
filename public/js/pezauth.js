var pezPortal = angular.module('pezPortal', [], function($interpolateProvider) {
      $interpolateProvider.startSymbol('{*{');
      $interpolateProvider.endSymbol('}*}');
  })
  .controller('PezPortalController', function($scope, $http, $timeout, $window) {
    $scope.hideCLIExample = true;
    var myData = {};
    var pauth = this;
    var restUriBase = "/v1/auth/api-key";
    var restOrgUriBase = "/v1/org/user";
    var meUri = "/me";
    var messaging = {
      "hasOrgBtn": "View Org Now",
      "createOrgBtn": "Create Your Org Now",
      "claimLease" : "Get One Now",
      "noApiKey": "You don't have a key yet",
      "loading": "Loading... Please Wait",
      "oktaSetup": "Get Okta Tile for HeritageCF",
      "invalidUser": "query failed. unable to find matching user guid."
    };

    $scope.claimButtonText = messaging.claimLease;
    $scope.hideClaimButton = false;

    // this will get dynamically populated by calling the local go pcfaas service
    // once created.
    //$scope.claimStatusText = "You have 1 week left on your lease of 2C.small.";
    //$scope.claimStatusText = "There are no available leases. The next one will be available in 3 days.";

    console.log('claim button text: ' + $scope.claimButtonText);
    console.log('hide claim button: ' + $scope.hideClaimButton);

    var urls = {
      "okta": "http://login.run.pez.pivotal.io/saml/login/alias/login.run.pez.pivotal.io?disco=true",
      "oktaHome": "https://pivotal.okta.com/app/UserHome"
    };

    $timeout(function () {
      callMeUsingVerb($http.get, meUri);
    }, 1);

    pauth.getRestUri = function() {
      return [restUriBase, $scope.myEmail].join("/");
    }

    pauth.getOrgRestUri = function() {
      return [restOrgUriBase, $scope.myEmail].join("/");
    }

    pauth.getorg = function() {
      console.log(pauth.getOrgRestUri());
      getOrgStatus(pauth.getOrgRestUri());
    };

    pauth.getInventory = function() {
      var uri = '/pcfaas/inventory';
      console.log('Getting inventory ' + uri );
      getInventoryList(uri);
    }

    pauth.createorg = function() {

      if ($scope.orgButtonText === messaging.createOrgBtn) {
        createOrg(pauth.getOrgRestUri());

      } else if ($scope.orgButtonText === messaging.hasOrgBtn) {
        $window.location.href = urls.okta;

      }  else if ($scope.orgButtonText === messaging.oktaSetup) {
        $window.location.href = urls.oktaHome;
      }
      $scope.orgButtonText = messaging.loading;
    };

    pauth.create = function() {
      callAPIUsingVerb($http.put, pauth.getRestUri());
    };

    pauth.remove = function() {
      callAPIUsingVerb($http.delete, pauth.getRestUri());
    };

    pauth.soonestexpiringinventoryitem = function() {
       return getSoonestExpiringInventoryItem($scope.inventoryItems);
    }

    pauth.firstavailableinventoryitem = function() {
      return getFirstAvailableInventoryItem($scope.inventoryItems);
    }

    pauth.myactiveleaseditem = function() {
      return getMyActiveLeasedInventoryItem($scope.inventoryItems);
    }

    pauth.setupLeasedItemState = function() {
      determineLeaseState();
    }

    pauth.leaseItem = function() {
      var firstAvailableInventoryItem = getFirstAvailableInventoryItem($scope.inventoryItems);
      return postInventoryItemRequest(firstAvailableInventoryItem)
    }

    function callMeUsingVerb(verbCaller, uri) {
      var responsePromise = verbCaller(uri);
      responsePromise.success(function(data, status, headers, config) {
          $scope.myName = data.Payload.displayName;
          $scope.myEmail = data.Payload.emails[0].value;
          $scope.displayName = $scope.myName ? $scope.myName : $scope.myEmail
          callAPIUsingVerb($http.get, pauth.getRestUri());
          pauth.getorg();
          pauth.getInventory();
      });
    }

    function getMyActiveLeasedInventoryItem(inventoryItems) {
      var result = inventoryItems.filter(function(el) {
          return el.currentLease != null && el.currentLease.userName === $scope.myEmail
      });
      if (result.length > 0) {
        return result[0];
      } else {
        return null;
      }
    }

    function getFirstAvailableInventoryItem(inventoryItems) {
      var result = inventoryItems.filter(function(el) {
        return el.status.toLowerCase() === "available"
      });
      if (result.length > 0) {
        return result[0];
      } else {
        return null;
      }
    }

    function getSoonestExpiringInventoryItem(inventoryItems) {
      var copiedItems = JSON.parse(JSON.stringify(inventoryItems))
      copiedItems = copiedItems.filter(function(el) {
        // could also get fancy here and return only inventory items of specific status
        return el.currentLease != null
      });
      copiedItems = copiedItems.sort(function(a, b){
          return a.currentLease.daysUntilExpires - b.currentLease.daysUntilExpires
      });
      if (copiedItems.length > 0) {
        return copiedItems[0];
      } else {
        return null;
      }
    }

    function createOrg(uri) {
      var responsePromise = $http.put(uri);
      responsePromise.success(function(data, status, headers, config) {
        console.log(data);
        $scope.orgButtonText = messaging.hasOrgBtn;
        $scope.hideCLIExample = false;
      });

      responsePromise.error(function(data, status, headers, config) {
          var forwardToOkta = false;

          if(status === 403) {
            console.log(data.ErrorMsg);

            if (messaging.invalidUser == data.ErrorMsg) {
              forwardToOkta = confirm("You have not set up your account in Okta yet. Please head over to Okta and click on the `PEZ HeritageCF` tile.");
            }

            if ( forwardToOkta === true) {
              $window.location.href = urls.oktaHome;

            } else {
              $scope.orgButtonText = messaging.oktaSetup;
            }
          }
      });
    }

    function getOrgStatus(uri) {
      var responsePromise = $http.get(uri);
      responsePromise.success(function(data, status, headers, config) {
        console.log(data);
        $scope.orgButtonText = messaging.hasOrgBtn;
        $scope.hideCLIExample = false;
      });

      responsePromise.error(function(data, status, headers, config) {

          if(status === 403) {
            $scope.orgButtonText = messaging.createOrgBtn;
            console.log(data.ErrorMsg);
          }
      });
    }

    function getInventoryList(uri) {
      var responsePromise = $http.get(uri);
      responsePromise.success(function(data, status, headers, config) {
        console.log(data);
        $scope.inventoryItems = data;
        determineLeaseState();
      });

      responsePromise.error(function(data, status, headers, config) {
        console.log(data.ErrorMsg);
      });
    }

    function postInventoryItemRequest(inventoryItem) {
      var uri = "/pcfaas/inventory/" + inventoryItem.id;
      var responsePromise = $http.post(uri);
      responsePromise.success(function(data, status, headers, config) {
        console.log(data);
        $scope.hideClaimButton = true;
        $scope.claimStatusText = "The lease on your inventory item " + inventoryItem.sku + " will expire in " + data.daysUntilExpires + " days.";
      });

      responsePromise.error(function(data, status, headers, config) {
        console.log(data.ErrorMsg);
        $scope.hideClaimButton = true;
        $scope.claimStatusText = "There was an error attempting to lease the " + inventoryItem.sku + " inventory item.";
      });
    }

    /*
     * Determines state of $scope variables: claimButtonText, claimStatusText, and hideClaimButton
     */
    function determineLeaseState() {
      var myLeasedItem = getMyActiveLeasedInventoryItem($scope.inventoryItems);
      var firstAvailableInventoryItem = getFirstAvailableInventoryItem($scope.inventoryItems);
      var soonestExpiringItem = getSoonestExpiringInventoryItem($scope.inventoryItems);

      if (myLeasedItem != null) {
        $scope.hideClaimButton = true;
        $scope.claimStatusText = "The lease on your inventory item " + myLeasedItem.sku + " will expire in " + myLeasedItem.currentLease.daysUntilExpires + " days.";
      } else if (firstAvailableInventoryItem != null) {
        $scope.hideClaimButton = false;
        $scope.claimButtonText = messaging.claimLease;
      } else if (soonestExpiringItem != null) { // someone else's lease is expiring soon.
        $scope.hideClaimButton = true;
        $scope.claimStatusText = "You will be able to claim a lease on a " + soonestExpiringItem.sku + " in " + soonestExpiringItem.currentLease.daysUntilExpires + " days.";
      } else {
        $scope.hideClaimButton = true;
        $scope.claimStatusText = "There are no inventory items, available or otherwise. PCFaaS may be down."
      }
    }

    function callAPIUsingVerb(verbCaller, uri) {
      var responsePromise = verbCaller(uri);

      responsePromise.success(function(data, status, headers, config) {
          $scope.myData = data;
          $scope.myApiKey = data.APIKey;
      });

      responsePromise.error(function(data, status, headers, config) {
        $scope.myApiKey = messaging.noApiKey;
        pauth.create();
      });
    }
  });
