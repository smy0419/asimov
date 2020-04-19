pragma solidity 0.4.25;

interface Warehouse {
	function init() external;
	function create(uint16 _category, string _name, bytes32 _key) external;
	function getSubmittedTemplatesCount(uint16 _category) external view returns (uint256);
	function getSubmittedTemplate(uint16 _category, uint256 index) external view returns (string, bytes32, uint, uint8, uint8, uint8, Status);
	function getApprovedTemplatesCount(uint16 _category) external view returns (uint256);
	function getApprovedTemplate(uint16 _category, uint256 index) external view returns (string, bytes32, uint, uint8, uint8, uint8, Status);
	function getTemplate(uint16 _category, string name) external view returns (string, bytes32, uint, uint8, uint8, uint8, Status);
	function disable(uint16 _category, string _name) external actionAuth(_category, _name);
	function enable(uint16 _category, string _name) external actionAuth(_category, _name);
	function nameExist(uint16 _category, string _name) external view returns (bool);
}

contract TestTemplateWarehouse {
	Warehouse internal warehouse;

	event LogResult(bool);

	function setUp() {
		warehouse = Warehouse(0x630000000000000000000000000000000000000067);
	}

	function testInit() {
		setUp();

		warehouse.init();
	}

	function testCreate() public {
		setUp();

		warehouse.create(1, "name1", 0x6300000000000000000000000000000000000000679999999999999999999999);
		/// failed; _category must bigger than 0
		warehouse.create(0, "name1", 0x6300000000000000000000000000000000000000679999999999999999999999);
		/// failed; template name can not bo empty
		warehouse.create(1, "", 0x6300000000000000000000000000000000000000679999999999999999999999);
		/// failed; file key's length must bigger than 0
		warehouse.create(1, "name2", 0x0000000000000000000000000000000000000000000000000000000000000000);
	}

	function testGetSubmittedTemplatesCount() public {
		setUp();

		warehouse.create(1, "name1", 0x6300000000000000000000000000000000000000679999999999999999999999);
		uint count = warehouse.getSubmittedTemplatesCount(1);
		if (1 == count) {
			emit LogResult(true);
		} else {
			emit LogResult(false);
		}
	}

	function testGetSubmittedTemplate() public {
		setUp();

		warehouse.create(1, "name1", 0x6300000000000000000000000000000000000000679999999999999999999999);

		string name;
		(name, , , , , , ) = warehouse.getSubmittedTemplate(1, 0);
		if (bytes(name).length > 0) {
			emit LogResult(true);
		} else {
			emit LogResult(false);
		}
	}

	function testGetApprovedTemplatesCount() public {
		setUp();

		warehouse.create(1, "name1", 0x6300000000000000000000000000000000000000679999999999999999999999);
		uint count = warehouse.getApprovedTemplatesCount(1, 0);
		if (1 == count) {
			emit LogResult(true);
		} else {
			emit LogResult(false);
		}
	}

	function testGetSubmittedTemplate() public {
		setUp();

		warehouse.create(1, "name1", 0x6300000000000000000000000000000000000000679999999999999999999999);

		string name;
		(name, , , , , , ) = warehouse.getApprovedTemplate(1, 0);
		if (bytes(name).length > 0) {
			emit LogResult(true);
		} else {
			emit LogResult(false);
		}
	}

	function testGetTemplate() public {
		setUp();

		warehouse.create(1, "name1", 0x6300000000000000000000000000000000000000679999999999999999999999);
		string name;
		(name, , , , , , ) = warehouse.getTemplate(1, "name1");
		if (bytes(name).length > 0) {
			emit LogResult(true);
		} else {
			emit LogResult(false);
		}
	}

	function testDisable() public {
		setUp();

		uint16 category = 1;
		string name = "name1";
		warehouse.create(category, name, 0x6300000000000000000000000000000000000000679999999999999999999999);
		string returnName;
		(returnName, , , , , , ) = warehouse.getTemplate(category, name);
		if (bytes(returnName).length > 0) {
			warehouse.disable(category, name);
			
			string returnName2;
			(returnName2, , , , , , ) = warehouse.getTemplate(category, name);
			if (bytes(returnName2).length == 0) {
				emit LogResult(true);
			} else {
				emit LogResult(false);
			}
		}
	}

	function testEnable() public {
		setUp();

		uint16 category = 1;
		string name = "name1";
		warehouse.create(category, name, 0x6300000000000000000000000000000000000000679999999999999999999999);
		string returnName;
		(returnName, , , , , , ) = warehouse.getTemplate(category, name);
		if (bytes(returnName).length > 0) {
			warehouse.disable(category, name);
			
			string returnName2;
			(returnName2, , , , , , ) = warehouse.getTemplate(category, name);
			if (bytes(returnName2).length == 0) {
				warehouse.enable(category, name);

				string returnName3;
				(returnName3, , , , , , ) = warehouse.getTemplate(category, name);
				if (bytes(returnName3).length > 0) {
					emit LogResult(true);
				} else {
					emit LogResult(false);
				}
			}
		}
	}

	function testNameExist() public {
		setUp();

		warehouse.create(1, "name1", 0x6300000000000000000000000000000000000000679999999999999999999999);
		bool result = warehouse.nameExist(1, "name1");
		if (result) {
			emit LogResult(true);
		} else {
			emit LogResult(false);
		}
	}

}




