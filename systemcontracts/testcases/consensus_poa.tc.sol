pragma solidity 0.4.25;
pragma experimental ABIEncoderV2;

interface POA {
	function init(address[] _admins, address[] _miners) external;
	function batchAddAdmins(address[] newAdmins) external;
	function batchRemoveAdmins(address[] oldAdmins) external;
	function batchAddValidators(address[] newValidators) external;
	function batchRemoveValidators(address[] oldValidataors) external;
	function getSignupValidators(string[] validatorsStr) external view returns(address[], uint[]);
	function getAdminsAndValidators() external view returns(address[], address[]);
}

contract TestConsensusPOA {
	POA internal poa;

	event LogResult(bool);
	event LogExistingMember(bool);

	function setUp() internal {
		poa = POA(0x63000000000000000000000000000000000000006c);
	}

	function testInit() public {
		setUp();

		address[] memory admins = new address[](2);
		admins[0] = 0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a;
		admins[1] = 0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad;

		address[] memory miners = new address[](2);
		miners[0] = 0x66a8c1def9e07589415ebf59b55f76b9d5115064c9;
		miners[1] = 0x662306d43258293846ab198810879e31feac759fac;
		poa.init(admins, miners);
	}

	function testBatchAddAdmins(address newAdmin) public {
		setUp();

		address[] memory admins;
		(admins, ) = poa.getAdminsAndValidators();

		bool existAdmin;
		for (uint i = 0; i < admins.length; i++) {
			if (newAdmin == admins[i]) {
				existAdmin = true;
			}
		}

		if (!existAdmin) {
			address[] memory newAdmins = new address[](1);
			newAdmins[0] = newAdmin;
			poa.batchAddAdmins(newAdmins);

			address[] memory returnAdmins;
			(returnAdmins, ) = poa.getAdminsAndValidators();
			for (uint j = 0; j < returnAdmins.length; j++) {
				if (newAdmin == returnAdmins[j]) {
					existAdmin = true;
				}
			}
			if (existAdmin) {
				emit LogResult(true);
			} else {
				emit LogResult(false);
			}
		} else {
			emit LogExistingMember(true);
		}
	}

	function testBatchRemoveAdmins(address oldAdmin) public {
		setUp();

		address[] memory admins;
		(admins, ) = poa.getAdminsAndValidators();

		bool existAdmin;
		for (uint i = 0; i < admins.length; i++) {
			if (oldAdmin == admins[i]) {
				existAdmin = true;
			}
		}
		if (!existAdmin) {
			emit LogExistingMember(true);
			return;
		}

		address[] memory oldAdmins = new address[](1);
		oldAdmins[0] = oldAdmin;
		poa.batchRemoveAdmins(oldAdmins);

		address[] memory returnAdmins;
		(returnAdmins, ) = poa.getAdminsAndValidators();
		for (uint j = 0; j < returnAdmins.length; j++) {
			if (oldAdmin == returnAdmins[j]) {
				existAdmin = true;
			}
		}
		if (existAdmin) {
			emit LogResult(true);
		} else {
			emit LogResult(false);
		}
	}

	function testBatchAddValidators(address newValidator) public {
		setUp();

		address[] memory validators;
		(validators, ) = poa.getAdminsAndValidators();

		bool existValidator;
		for (uint i = 0; i < validators.length; i++) {
			if (newValidator == validators[i]) {
				existValidator = true;
			}
		}

		if (!existValidator) {
			address[] memory newValidators = new address[](1);
			newValidators[0] = newValidator;
			poa.batchAddValidators(newValidators);

			address[] memory returnValidators;
			(returnValidators, ) = poa.getAdminsAndValidators();
			for (uint j = 0; j < returnValidators.length; j++) {
				if (newValidator == returnValidators[j]) {
					existValidator = true;
				}
			}
			if (existValidator) {
				emit LogResult(true);
			} else {
				emit LogResult(false);
			}
		} else {
			emit LogExistingMember(true);
		}
	}

	function testBatchRemoveValidators(address oldValidator) public {
		setUp();

		address[] memory validators;
		(validators, ) = poa.getAdminsAndValidators();

		bool existValidator;
		for (uint i = 0; i < validators.length; i++) {
			if (oldValidator == validators[i]) {
				existValidator = true;
			}
		}
		if (!existValidator) {
			emit LogExistingMember(true);
			return;
		}

		address[] memory oldValidators = new address[](1);
		oldValidators[0] = oldValidator;
		poa.batchRemoveValidators(oldValidators);

		address[] memory returnValidators;
		(returnValidators, ) = poa.getAdminsAndValidators();
		for (uint j = 0; j < returnValidators.length; j++) {
			if (oldValidator == returnValidators[j]) {
				existValidator = true;
			}
		}
		if (existValidator) {
			emit LogResult(true);
		} else {
			emit LogResult(false);
		}
	}

	function testGetSignupValidators() public {
		setUp();

		string[] memory mappingStrs = new string[](0);

		address[] memory validators;
		uint[] memory intArray;
		(validators, intArray) = poa.getSignupValidators(mappingStrs);
		if (validators.length == intArray.length) {
			emit LogResult(true);
		} else {
			emit LogResult(false);
		}
	}

	function testGetAdminsAndValidators() public {
		setUp();

		address[] memory admins;
		address[] memory validators;
		(admins, validators) = poa.getAdminsAndValidators();
	}

}


