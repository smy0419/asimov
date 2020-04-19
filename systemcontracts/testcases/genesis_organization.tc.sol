pragma solidity 0.4.25;

import "../utils/test_tool.sol";
import "../utils/execution.sol";

interface GenesisOrganization {
	function init(address[] citizens) external;
	function startProposal(uint proposalType, address targetAddress, uint amount, uint asset) external returns(uint);
	function vote(uint proposalId, bool decision) external;
	function testGetProposalDetail(uint proposalId) external view returns(address, address, ProposalStatus, address[], address[], address[], ProposalType, uint, uint);
}

interface Registry {
	function newAsset(string name, string symbol, string description, uint32 assetType, uint32 assetIndex, uint amountOrVoucherId) external;
}

contract TestGenesisOrganization is TestTool, Execution {
	GenesisOrganization internal organization;
	Registry internal registry;

	event LogResult(bool);

	function setUp() internal {
		organization = GenesisOrganization(0x630000000000000000000000000000000000000064);
	}

	function createAsset() internal {
		registry = Registry(0x630000000000000000000000000000000000000065);
		flow.createAsset(0, 1, 1000000000);
        registry.newAsset("test coin", "coin", "this is test coin", 0, 1, 1000000000);
	}

	function testFallBack() public {
		setUp();
		createAsset();

        address(organization).transfer(7777777, 4294967297);
	}

	function testStartProposal() public {
		setUp();

		bool result = execute(organization, "startProposal(uint,address,uint,uint)", abi.encode(0,0x66ad16b77bffe22b728374e717867a6b1042418c78,0,0));
		assertTrue(result);

		result = execute(organization, "startProposal(uint,address,uint,uint)", abi.encode(0,0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a,0,0));
		assertTrue(!result);

		result = execute(organization, "startProposal(uint,address,uint,uint)", abi.encode(1,0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a,0,0));
		assertTrue(result);

		createAsset();
		address(organization).transfer(8888888, 4294967297);
		organization.startProposal(2, 0x0, 111, 4294967297);
		organization.startProposal(2, 0x0, 222, 4294967297);

		/// above will fail
		organization.startProposal(0, 0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a, 0, 0);
		organization.startProposal(1, 0x665c48f910524a9395e067d489a9c5a9041bc0c35b, 0, 0);
		organization.startProposal(2, 0x0, 999999999999, 4294967297);
		organization.startProposal(2, 0x0, 999, 4294967296);
	}

	function testStartProposal2() public {
		setUp();

		uint id = organization.startProposal(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff, 0x66ad16b77bffe22b728374e717867a6b1042418c78, 0, 0);
		if (id > 0) {
			emit LogResult(false);
		}
	}

	function testVote() public {
		setUp();

		uint proposalId = organization.startProposal(0, 0x66ad16b77bffe22b728374e717867a6b1042418c78, 0, 0);
		organization.vote(proposalId, true);

		uint proposalId2 = organization.startProposal(1, 0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a, 0, 0);
		organization.vote(proposalId2, false);

		uint proposalId3 = organization.startProposal(0, 0x66ad16b77bffe22b728374e717867a6b1042418c78, 0, 0);
		organization.vote(proposalId3+1, true);
	}

}