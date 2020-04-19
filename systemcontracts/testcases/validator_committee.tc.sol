pragma solidity 0.4.25;

interface Committee {
	function init(address[] _validators) external;
	function signup() external;
	function startNewRound(address[] validatorAddresses, uint[] plannedBlocks, uint[] actualBlocks) external;
	function startProposal(uint proposalType, uint asset) external returns(uint);
	function vote(uint proposalId, bool decision) external;
	function getAssetFeeList() external view returns(uint[], uint[]);
	function testGetProposalDetail(uint proposalId) external view returns(address, ProposalStatus, address[], address[], address[]);
}


contract TestCommittee {
	Committee internal committee;

	function setUp() public {
		committee = Committee(0x63000000000000000000000000000000000000006b);
	}

	function testInit() public {
		setUp();

		address[] memory validators = new address[](5);
		validators[0] = 0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a;
		validators[1] = 0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad;
		validators[2] = 0x66a8c1def9e07589415ebf59b55f76b9d5115064c9;
		validators[3] = 0x662306d43258293846ab198810879e31feac759fac;
		validators[4] = 0x6697a7f6e03dcff5cd9cbf99b5a77e7b31dd07153c;
		committee.init(validators);
	}

	function testSignup() public {
		setUp();

		committee.signup();
	}

	function testStartNewRound() public {
		setUp();

		address[] memory validators = new address[](5);
		validators[0] = 0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a;
		validators[1] = 0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad;
		validators[2] = 0x66a8c1def9e07589415ebf59b55f76b9d5115064c9;
		validators[3] = 0x662306d43258293846ab198810879e31feac759fac;
		validators[4] = 0x6697a7f6e03dcff5cd9cbf99b5a77e7b31dd07153c;

		uint[] memory planned = new uint[](5);
		planned[0] = 100;
		planned[1] = 100;
		planned[2] = 100;
		planned[3] = 100;
		planned[4] = 100;

		uint[] memory actual = new uint[](5);
		actual[0] = 90;
		actual[1] = 80;
		actual[2] = 70;
		actual[3] = 60;
		actual[4] = 50;

		 committee.startNewRound(validators, planned, actual);
	}

	function testStartProposal() public {
		setUp();

		uint proposalId = committee.startProposal(0, 4294967297);
		address proposer;
		(proposer, , , , ) = committee.testGetProposalDetail(proposalId);
		if (0x0 != proposer) {
			emit LogResult(true);
		} else {
			emit LogResult(false);
		}

		/// failed; proposalType must be 0
		committee.startProposal(1, 4294967297);
	}

	function testVote() public {
		setUp();

		uint proposalId = committee.startProposal(0, 4294967297);
		committee.vote(proposalId, true);

		uint proposalId2 = committee.startProposal(0, 4294967298);
		committee.vote(proposalId2, false);

		/// failed; can not vote twice
		uint proposalId3 = committee.startProposal(0, 4294967299);
		committee.vote(proposalId3, true);
		committee.vote(proposalId3, false);
	}

	function testGetAssetFeeList() public {
		setUp();

		committee.getAssetFeeList();
	}

}




