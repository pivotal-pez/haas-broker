describe('PezPortalController', function() {
  beforeEach(module('pezPortal'));

  var $controller;
  var $httpBackend;
  var testEmailAddy = 'test@pivotal.io';
  var testName = "Testy Larue";
  var testAPIKey = "12345";

  function createLease(days) {
    return {"daysUntilExpires": days, "userName": "foo@bar.com"}
  }

  function createSoonExpiringInventoryItems() {
    return [
      {"sku":"2C.small","tier":"2","offeringType":"C","size":"small","status":"available","id":"abc123guid"},
      {"sku":"2C.small","tier":"2","offeringType":"C","size":"small","status":"leased","id":"abc234guid","currentLease" : createLease(3)},
      {"sku":"2C.small","tier":"2","offeringType":"C","size":"small","status":"leased","id":"abc456guid","currentLease" : createLease(7)}
    ];
  }

  function createInventoryItemsWithMixedAvailability() {
    return [
      {"sku":"2C.small","tier":"2","offeringType":"C","size":"small","status":"leased","id":"abc123guid","currentLease" : createLease(9)},
      {"sku":"2C.small","tier":"2","offeringType":"C","size":"small","status":"available","id":"abc234guid"},
      {"sku":"2C.small","tier":"2","offeringType":"C","size":"small","status":"available","id":"abc456guid"}
    ];
  }

  function createNoAvailableInventoryItems() {
    return [
      {"sku":"2C.small","tier":"2","offeringType":"C","size":"small","status":"leased","id":"abc123guid","currentLease" : createLease(9)},
      {"sku":"2C.small","tier":"2","offeringType":"C","size":"small","status":"leased","id":"abc234guid","currentLease" : createLease(5)}
    ];
  }

  function createInventoryItemsWithNoPendingExpirations() {
    return [
      {"sku":"2C.small","tier":"2","offeringType":"C","size":"small","status":"available","id":"abc123guid"},
      {"sku":"2C.small","tier":"2","offeringType":"C","size":"small","status":"available","id":"abc234guid"},
      {"sku":"2C.small","tier":"2","offeringType":"C","size":"small","status":"available","id":"abc456guid"}
    ];
  }

  beforeEach(inject(function($injector){
    $httpBackend = $injector.get('$httpBackend');
    $controller = $injector.get('$controller');
    $window = $injector.get('$window');
  }));

  describe('$scope.myName & myEmail', function() {
    it('should be initialized as undefined', function() {
      var $scope = {};
      var controller = $controller('PezPortalController', { $scope: $scope });
      expect($scope.myName).toEqual(undefined);
      expect($scope.myEmail).toEqual(undefined);
    });
  });

  describe('$scope.myName & myEmail', function() {
    it('should allow initialization', function() {
      var $scope = {"myEmail": testEmailAddy, "myName": testName};
      var controller = $controller('PezPortalController', { $scope: $scope });
      expect($scope.myName).toEqual(testName);
      expect($scope.myEmail).toEqual(testEmailAddy);
    });
  });

  describe('$scope.hideCLIExample', function() {
    it('should initialize to true', function() {
      var $scope = {};
      var controller = $controller('PezPortalController', { $scope: $scope });
      expect($scope.hideCLIExample).toEqual(true);
    });
  });

  describe('$scope.hideClaimButton', function() {
    it('should initialize to false', function() {
      var $scope = {};
      var controller = $controller('PezPortalController', { $scope: $scope });
      expect($scope.hideClaimButton).toEqual(false);
    });
  });

  describe('pezauth.soonestexpiringinventoryitem', function() {
    it('should return the right item among multiple expiring items', function() {
        var $scope = {};
        $scope.inventoryItems = createSoonExpiringInventoryItems();
        var originalLength = $scope.inventoryItems.length;
        var controller = $controller('PezPortalController', { $scope: $scope });
        var item = controller.soonestexpiringinventoryitem();
        expect(item.id).toEqual("abc234guid")
        expect(originalLength).toEqual($scope.inventoryItems.length);
    });

    it('should return nil when no leases are outstanding', function() {
      var $scope = {};
      $scope.inventoryItems = createInventoryItemsWithNoPendingExpirations();
      var originalLength = $scope.inventoryItems.length;
      var controller = $controller('PezPortalController', { $scope: $scope });
      var item = controller.soonestexpiringinventoryitem();
      expect(item).toBe(null);
      expect(originalLength).toEqual($scope.inventoryItems.length);
    });
  });

  describe('pezauth.myactiveleaseditem', function() {
    it('should return the right inventory item based on email address', function() {
        var $scope = {};
        $scope.inventoryItems = createSoonExpiringInventoryItems();
        var originalLength = $scope.inventoryItems.length;
        $scope.myEmail = "foo@bar.com";
        var controller = $controller('PezPortalController', { $scope: $scope });
        var item = controller.myactiveleaseditem();
        expect(item.id).toEqual("abc234guid");
        expect(originalLength).toEqual($scope.inventoryItems.length);
    });

    it('should return null when no inventory lease matches', function() {
      var $scope = {};
      $scope.inventoryItems = createSoonExpiringInventoryItems();
      var originalLength = $scope.inventoryItems.length;
      $scope.myEmail = "definitely@doesnt.exist.com";
      var controller = $controller('PezPortalController', { $scope: $scope });
      var item = controller.myactiveleaseditem();
      expect(item).toBe(null);
      expect(originalLength).toEqual($scope.inventoryItems.length);
    });
  })

  describe('pezauth.firstavailableinventoryitem', function() {
    it('should return the first item with status available', function() {
      var $scope = {};
      $scope.inventoryItems = createInventoryItemsWithMixedAvailability();
      var originalLength = $scope.inventoryItems.length;
      var controller = $controller('PezPortalController', { $scope: $scope });
      var item = controller.firstavailableinventoryitem();
      expect(item.id).toEqual("abc234guid");
      expect(originalLength).toEqual($scope.inventoryItems.length); // assert we didn't modify the raw items list
    });

    it('should return null when no available items', function() {
      var $scope = {};
      $scope.inventoryItems = createNoAvailableInventoryItems();
      var originalLength = $scope.inventoryItems.length;
      var controller = $controller('PezPortalController', { $scope: $scope });
      var item = controller.firstavailableinventoryitem();
      expect(item).toBe(null);
      expect(originalLength).toEqual($scope.inventoryItems.length); // assert we didn't modify the raw items list
    });
  })

  describe('Determine inventory state', function() {
    it('should show first available item when one is available', function() {
      var $scope = {};
      $scope.inventoryItems = createSoonExpiringInventoryItems();
      var controller = $controller('PezPortalController', { $scope: $scope });
      controller.setupLeasedItemState();
      expect($scope.hideClaimButton).toBe(false);
      expect($scope.claimButtonText).toBe("Get One Now");
    });
    it('should show current user expiring lease', function() {
      var $scope = {};
      $scope.myEmail = "foo@bar.com";
      $scope.inventoryItems = createSoonExpiringInventoryItems();
      var controller = $controller('PezPortalController', { $scope: $scope });
      controller.setupLeasedItemState();
      expect($scope.hideClaimButton).toBe(true);
      expect($scope.claimStatusText).toBe("The lease on your inventory item 2C.small will expire in 3 days.");
    });
    it('should show when next available lease is coming up when user has none and none are available', function() {
      var $scope = {};
      $scope.inventoryItems = createNoAvailableInventoryItems();
      var controller = $controller('PezPortalController', { $scope: $scope });
      controller.setupLeasedItemState();
      expect($scope.hideClaimButton).toBe(true);
      expect($scope.claimStatusText).toBe("You will be able to claim a lease on a 2C.small in 5 days.");
    });
  });

  describe('pezportal.Lease inventory item', function() {
    it('should return a lease object when successful', function() {
      var inventoryItemID = "abc123guid"; // this is the ID of the "first available" item.
      $httpBackend.when('POST', ['/pcfaas/inventory', inventoryItemID].join('/')).respond(
        {"daysUntilExpires": 14, "userName" : "bob"}
      );
      var $scope = {};
      $scope.inventoryItems = createSoonExpiringInventoryItems();
      var controller = $controller('PezPortalController', { $scope: $scope });
      controller.leaseItem();
      $httpBackend.flush();

      expect($scope.hideClaimButton).toBe(true);
      expect($scope.claimStatusText).toBe("The lease on your inventory item 2C.small will expire in 14 days.");
    });

    it('should expose cause for failure when lease fails', function() {
      var inventoryItemID = "abc123guid"; // this is the ID of the "first available" item.
      $httpBackend.when('POST', ['/pcfaas/inventory', inventoryItemID].join('/')).respond(500,'');
      var $scope = {};
      $scope.inventoryItems = createSoonExpiringInventoryItems();
      var controller = $controller('PezPortalController', { $scope: $scope });
      controller.leaseItem();
      $httpBackend.flush();

      expect($scope.hideClaimButton).toBe(true);
      expect($scope.claimStatusText).toBe("There was an error attempting to lease the 2C.small inventory item.");
    });
  })

   describe('pezportal.getOrgRestUri', function() {
    it('should combine the ORG API base with email into a path', function() {
      var $scope = {"myEmail": testEmailAddy};
      var controller = $controller('PezPortalController', { $scope: $scope });
      expect(controller.getOrgRestUri()).toEqual('/v1/org/user/' + testEmailAddy);
    });
  });

  describe('pezportal.getRestUri', function() {
    it('should combine the API base with email', function() {
      var $scope = {"myEmail": testEmailAddy};
      var controller = $controller('PezPortalController', { $scope: $scope });
      expect(controller.getRestUri()).toEqual('/v1/auth/api-key/' + testEmailAddy);
    });
  });

  describe('pezportal.create', function() {
    it('should create an API Key', function() {
      $httpBackend.when('PUT', ['/v1/auth/api-key', testEmailAddy].join('/')).respond({"APIKey": testAPIKey});

      var $scope = {"myEmail": testEmailAddy};
      var controller = $controller('PezPortalController', { $scope: $scope });
      controller.create();
      $httpBackend.flush();

      expect($scope.myData.APIKey).toBe(testAPIKey);
      });
  });

  describe('pezportal.createorg', function() {
    it('should create an createorg when orgButtonText is "Create Your Org Now"', function() {
      $httpBackend.when('PUT', ['/v1/org/user', testEmailAddy].join('/')).respond(201);

      var $scope = {"myEmail": testEmailAddy, "orgButtonText": "Create Your Org Now"};
      var controller = $controller('PezPortalController', { $scope: $scope });
      controller.createorg();
      $httpBackend.flush();

      expect($scope.orgButtonText).toBe("View Org Now");
    });
  });

  describe('pezportal.createorg', function() {
    it('should not create an org when a user can\'t be found', function() {
      $httpBackend.when('PUT', ['/v1/org/user', testEmailAddy].join('/')).respond(403, {"ErrorMsg": "query failed. unable to find matching user guid."});

      var $scope = {"myEmail": testEmailAddy, "orgButtonText": "Create Your Org Now"};
      var controller = $controller('PezPortalController', { $scope: $scope });
      controller.createorg();
      $httpBackend.flush();

      expect($scope.orgButtonText).toBe("Get Okta Tile for HeritageCF");
    });
  });
});
