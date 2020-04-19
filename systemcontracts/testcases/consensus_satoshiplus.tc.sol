pragma solidity 0.4.25;
pragma experimental ABIEncoderV2;

import "../utils/test_tool.sol";
import "../utils/execution.sol";

interface ConsensusSatoshiplus {
	function init() external;
	function signup() external;
	function getSignupValidators(string[] str) external returns(address[], uint[]);
    function getValidatorBlockInfo(address validatorAddress) external view returns(uint, uint, uint);
    function testGetMappingValidators() external view returns(address[]);
    function testGetSignupValidator(address validatorAddress) external view returns(bool);
}

contract TestConsensusSatoshiplus is TestTool, Execution {
	ConsensusSatoshiplus internal consensus;

	event LogResult(bool);

	function setUp() internal {
		consensus = ConsensusSatoshiplus(0x63000000000000000000000000000000000000006a);
	}

	/// signup will fail,contract should get mapping address first
	function testSignup() public {
		setUp();
		consensus.signup();

		bool result = consensus.testGetSignupValidator(msg.sender);
		assertTrue(result);
	}

	function testGetSignupValidators() public {
		setUp();

		consensus.signup();

		string[] memory str = new string[](2);
		str[0] = "aaa";
		str[1] = "bbb";

		address[] memory validators;
		uint[] memory rounds;
		(validators, rounds) = consensus.getSignupValidators(str);

		assertEq(validators[0], 0x0);
		assertEq(validators[0], validators[1]);
		assertEq(rounds[0], 0xffffffffffffffff);
		assertEq(rounds[0], rounds[1]);
	}

	function testUpdateValidatorsBlockInfo() public {
		setUp();
		address[] memory validators = new address[](2);
		validators[0] = 0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a;
		validators[1] = 0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad;

		uint[] memory planned = new uint[](2);
		planned[0] = 100;
		planned[1] = 100;

		uint[] memory actual = new uint[](2);
		actual[0] = 80;
		actual[1] = 90;

		string[] memory str = new string[](2);
		str[0] = "btc miner1";
		str[1] = "btc miner2";

		address[] memory mappingAddresses = new address[](2);
		mappingAddresses[0] = 0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a;
		mappingAddresses[1] = 0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad;

		string[] memory domains = new string[](2);
		domains[0] = "mining pool1";
		domains[1] = "mining pool2";

		bool result = execute(consensus, "updateValidatorsBlockInfo(address[],uint[],uint[],string[],address[],uint,string[])", abi.encode(validators,planned,actual,str,mappingAddresses,1,domains));
		assertTrue(result);
		uint plannedBlocks;
		uint actualBlocks;
		uint efficiency;
		(plannedBlocks, actualBlocks, efficiency) = consensus.getValidatorBlockInfo(validators[0]);
		assertEq(plannedBlocks, 100);
		assertEq(actualBlocks, 80);
		assertEq(efficiency, 80);

		/// 出块人参数少一个
		validators = new address[](1);
		validators[0] = 0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad;
		result = execute(consensus, "updateValidatorsBlockInfo(address[],uint[],uint[],string[],address[],uint,string[])", abi.encode(validators,planned,actual,str,mappingAddresses,1,domains));
		assertTrue(!result);

		/// 验证者数量少一个
		str = new string[](1);
		str[0] = "btc miner1";
		result = execute(consensus, "updateValidatorsBlockInfo(address[],uint[],uint[],string[],address[],uint,string[])", abi.encode(validators,planned,actual,str,mappingAddresses,1,domains));
		assertTrue(!result);

		/// 验证者数量少两个
		str = new string[](0);
		result = execute(consensus, "updateValidatorsBlockInfo(address[],uint[],uint[],string[],address[],uint,string[])", abi.encode(validators,planned,actual,str,mappingAddresses,1,domains));
		assertTrue(!result);

		/// 都是空的映射key，不同的映射地址，会有怎样的结果
		str = new string[](1);
		mappingAddresses = new address[](1);
		mappingAddresses[0] = 0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a;
		result = execute(consensus, "updateValidatorsBlockInfo(address[],uint[],uint[],string[],address[],uint,string[])", abi.encode(validators,planned,actual,str,mappingAddresses,1,domains));
		assertTrue(result);
		mappingAddresses[0] = 0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad;
		result = execute(consensus, "updateValidatorsBlockInfo(address[],uint[],uint[],string[],address[],uint,string[])", abi.encode(validators,planned,actual,str,mappingAddresses,1,domains));
		assertTrue(result);
		address[] memory returnValidators = consensus.testGetMappingValidators(str);
		assertEq(returnValidators[0], 0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad);
	}

	function testPressureUpdateValidatorsBlockInfo() public {
		setUp();
		address[] memory validators = new address[](120);
		validators[0] = 0x66d77a88053271d9e1639a25585dd10795afd2c3b7;
		validators[1] = 0x66da47d716ad3f6aaf5b5dd493da0ec013f50300d3;
		validators[2] = 0x667566d4aef7b44b05753aa6888ee1ad0f675d4146;
		validators[3] = 0x66be27459adbe053853d0fbb37c0f6a16b31c8f9fc;
		validators[4] = 0x66ae49d4f28b7a8b28fc900e62d2a4493098895a25;
		validators[5] = 0x666346bdacc274ec340f8943e8b47873bc1b89ff1b;
		validators[6] = 0x66361179a1ff3ffffcf7765d6f378cdfe440fb54c9;
		validators[7] = 0x6690793f5686c9efce4e024fee624f942c9412e05a;
		validators[8] = 0x6682fffd2e0c322cf9ff41cef6496f4ccc63971964;
		validators[9] = 0x666a375d2083e9cd93f286f558095d2677cfee7b46;
		validators[10] = 0x660480d495eacd4d4c713b403c095e089c23e9a10c;
		validators[11] = 0x66f110bb0a7cc08d62cd4c950994a78badbdbef38d;
		validators[12] = 0x6632fda273f14d9f5d1f8f10844c35d77e40be209d;
		validators[13] = 0x666a217405ce79ac1d50f5ea6f203c6a00d72b05c6;
		validators[14] = 0x66dc6a7c6e165ea7f71bc9adade0fafaf3b0e9064b;
		validators[15] = 0x669dd9b57774039a76079846dde072347602079ff9;
		validators[16] = 0x666cbf10d8ad236ba484b88bb9e168cf7498eb7a9f;
		validators[17] = 0x66f7f8f67b788cb14f08cbdbe02f65875f08dec8b0;
		validators[18] = 0x6690b6faf8b3cab9c16ac777a05ddb198f515783aa;
		validators[19] = 0x660e4c5e9f03328a95c2257348bcd2b96b07ae721b;
		validators[20] = 0x66ebd6e8670cffd3cad193a162945d4ef4c5033eb1;
		validators[21] = 0x66b29727f99199547f8827fac2e27fc1060de4fe32;
		validators[22] = 0x66b0a8b65f1db750b0938ba67975aedc330f638168;
		validators[23] = 0x6689ee0190cc61da9b6d394b3b87651e12eb625d85;
		validators[24] = 0x6688e50d4d5ada5c464f9f3e97ee09c775d5022812;
		validators[25] = 0x662ce04109187f03fcfa45651586320316ee8dd73e;
		validators[26] = 0x66a35bde893f1214201557bf2d5cc2480f2d65001a;
		validators[27] = 0x6639178ce239ecabb9d4edd8c18dd82b7815a5ae7d;
		validators[28] = 0x6681945a5ae337ed1b750298c510e916f095548011;
		validators[29] = 0x66bbf86d966a65d3ba6e2103b0d003ec2594d255fb;
		validators[30] = 0x6657e2b890339667785f1cfd1730f0cf4dbf4e2173;
		validators[31] = 0x66e6acb97e3ca6147cf7ddbc17ead96833444ab687;
		validators[32] = 0x664625530cee4355aa97b21a7b528385be0d1dbd09;
		validators[33] = 0x6658daa852911251dca0b0f308860bd2d33f929d8b;
		validators[34] = 0x66e14a19a82b99444eb68c621bf91eae353b32a63c;
		validators[35] = 0x6616b1070418b3d3cce1fd800e2261c6a7ee01eb9e;
		validators[36] = 0x66aad2eb0f4e5ca84a19eafac0dffc3a701f091737;
		validators[37] = 0x663327bf38ed4bd8911c4d0a3965d6414b443195ec;
		validators[38] = 0x66442986531a79ab300ea4035f6661825076604c13;
		validators[39] = 0x66eb694db040c0afd2e500313167e638ac71f8de82;
		validators[40] = 0x66aa9a3ae9eab6af4b8986cc322d9758234ff4e2e4;
		validators[41] = 0x6687ca76e938c06b660cf9d6f48adfd53db3bcb1e7;
		validators[42] = 0x666a33685cf9a7884f02b91afdbf30a97c123db8dc;
		validators[43] = 0x66c0f4d35218d930b82b85c09154e1e0b523e4259d;
		validators[44] = 0x662a04076848150e6e6671d436bbedf29aa949d9a7;
		validators[45] = 0x668f98db4d173a9a83bc899be348886cd1510fb843;
		validators[46] = 0x66a526e8c8593f9b9d027ad175f79572fe3834b398;
		validators[47] = 0x66047bac4389ea149a87dd6eda3002bef683ea5a31;
		validators[48] = 0x66b2c836b753cb906dde1d02f07f07c60eea94ff7f;
		validators[49] = 0x6638c62bc831211628b386f03b7169aeb5048d480c;
		validators[50] = 0x66b36e89e52ae2559d41cba19167e825947a3ff1da;
		validators[51] = 0x66fff93a61ee68229959d4603daa722ea55a0c2d2f;
		validators[52] = 0x66ae47551e6766e21a6b682ccd4329cbb4a1374e2b;
		validators[53] = 0x661e6d86e92bd2310e21833a1d58c97f281e09545f;
		validators[54] = 0x661a13ccb9e7e2b2ec581ad8ae83548ffacfa932fa;
		validators[55] = 0x661ca2fe435c240225ec967bb65a0b187361b92bb4;
		validators[56] = 0x6616ccc91f27ce26fc2e13d61c5869bfd06c0c629a;
		validators[57] = 0x663593c4b045c8c757bc3db294e923b94b9f7455a2;
		validators[58] = 0x66187da8cfdee52c7b19dff16aac1ce26085533030;
		validators[59] = 0x66bd589b6aafd37cb8f1e1de589cb79b001cf2a40c;
		validators[60] = 0x661c1a0ba435184779f446ed39822fa9f2a88ff040;
		validators[61] = 0x669d84a0826d288e585f70909d353421fecf003462;
		validators[62] = 0x665f809bb9015abee935c4c9d81342d2d2a9b8ab95;
		validators[63] = 0x66c7e89305e186125659942d6a561487390886643b;
		validators[64] = 0x665d9de005f73f4bbcc6a55a84d5119576b4cf8281;
		validators[65] = 0x663e9abbd9219bc5142b88b8dcee554fab3fd7a3ab;
		validators[66] = 0x6652f60a5830d1069058c0fc13568ed45294a8a1ea;
		validators[67] = 0x6601653342f55c20cad6ad96f40c14808705ea7e89;
		validators[68] = 0x6629a11b796b695e47d054c32bbbc32831125016ae;
		validators[69] = 0x66546c7e78da17136660a65cf250e1c5e9d2de00bb;
		validators[70] = 0x66d8f069bf3abdcd0c06abe458b2cc7d009e4c3608;
		validators[71] = 0x6682161c972116103e884f01976a7e940b9e15d623;
		validators[72] = 0x663a735ca225a3c1980f095f5f0e7607c29f199be6;
		validators[73] = 0x664adbc953bf91973bcc4ffb4be6716d25fd270449;
		validators[74] = 0x66c02fd34a529dd0d6f3eba6836a393d67281c2fc7;
		validators[75] = 0x6673cfdc1c7097cfeca8f1117528ba7e3c887f5076;
		validators[76] = 0x66b2e08a1d67e90b46f4fcf1a0e4890a85aeadddc8;
		validators[77] = 0x665d91287792316fc2218866eb162d6d00c5699407;
		validators[78] = 0x66da640ca837ffafa75a4c700b74d5a9aba13a8f29;
		validators[79] = 0x66d962f2bff531142889b125ef8fd43a1836fc9f9f;
		validators[80] = 0x66f4f4e24bffa164ed3805a5e337632de56ed99645;
		validators[81] = 0x66ecf0c0d4ee2c2c279dd1443bb8e6f724eb1888f3;
		validators[82] = 0x66e6bbfc17e398b946658394aa82b19130bb3e81d5;
		validators[83] = 0x66667b77ff427317a826310c3c52f736e81cad0b67;
		validators[84] = 0x66bcc7c06896c1cef1dfb22ed87a9eb0ec410db435;
		validators[85] = 0x66ab40a203c5063f4daaff1099b7e52aa76a19f5ac;
		validators[86] = 0x66f64153ec576569c84dfb80c65f4975cd0a47ceeb;
		validators[87] = 0x66ee635a0a3dafb31029afc233d907bd30ac0ed280;
		validators[88] = 0x669c23e3bff3bc81ecacbd5d224b505c6a8c28d5b5;
		validators[89] = 0x6675808c2a2a45c2aea94ec7130d9aaca5323d3dcd;
		validators[90] = 0x6685d72b21ba9d1e2c70be83052eb718f6f4de3370;
		validators[91] = 0x666c51bed6a55398ad5eb3c61e2e66c3c9381e031f;
		validators[92] = 0x664ee4d1d4e2fe722c4e50cfdf774536c01d0e9166;
		validators[93] = 0x669c36c7e36eeacce77569c3eeb8ea561b90c938ff;
		validators[94] = 0x662c31fb6061b59ee9f0e64fedee948ba2cbf540e1;
		validators[95] = 0x66515b0aa4eca369863579099fef78614ab8ff41ca;
		validators[96] = 0x66d3c6e124d6feef6d1540daa74198169f11c248e1;
		validators[97] = 0x66a7383846e892554acebfeb5a460f5960d0fed500;
		validators[98] = 0x66418a9fd322b319b1e34b73398e42aa2c72cd3d4d;
		validators[99] = 0x668a67319e293669c69a0af4a19a0bb4778e23be2f;
		validators[100] = 0x66255f181037f4a4646cff846421e54b43f314144a;
		validators[101] = 0x66617a89a6ed3b6902806849e92fdef68348c24f17;
		validators[102] = 0x66e20e3b4f73fce0ee2f58a5e0a4d5f3d273771858;
		validators[103] = 0x6663e56f4e5e1dd38be35e00efb55896894e6c744c;
		validators[104] = 0x66723f72e2d254d0b3b7d75e5148082910a23e3a95;
		validators[105] = 0x66e7dae7c5824b3ee7c083c69265213eee23b49e37;
		validators[106] = 0x66ec6df181530303aea17d6cbda85cce03e34be8c4;
		validators[107] = 0x66bbdec02b194254043fcbbede19163577b3bd02b0;
		validators[108] = 0x662defb867b4fd9ed454b0cd6bcc2c626299b5d04f;
		validators[109] = 0x6618b2812bf29c6afbd8a24c4a1bedf80f62aee4d1;
		validators[110] = 0x66b90c46986f41eaf64603eb7dcfcfb10e419a142f;
		validators[111] = 0x6669fc8b097a6e48a76481318b247dcb3ca2b2b31c;
		validators[112] = 0x66ac816861fd8e0c82ef11526fb5aa051ebb29a90b;
		validators[113] = 0x66e47cb74127c1844817024b4269fa7e4c1e382cde;
		validators[114] = 0x66b3fa69b576f46115a822c24fd9ee7caf6534773e;
		validators[115] = 0x6657d3873d03cef81d53de04075c0be7df4f4801fb;
		validators[116] = 0x66d3e5333b843ddfdcefc4b669b40d576e40106a6f;
		validators[117] = 0x660a4a23095995d77618af68310e861c22153538a0;
		validators[118] = 0x669db9f5bbedc10d8006227d3d1440f5e3a5b9bd1c;
		validators[119] = 0x66979ed5a14a05aaf1b23f916430b6537d109534ce;
		validators[120] = 0x66f1dec1bc8c452a00817af25fee7ffae5e1a62ec1;

		uint[] memory planned = new uint[](120);
		planned[0] = 100;
		planned[1] = 100;
		planned[2] = 100;
		planned[3] = 100;
		planned[4] = 100;
		planned[5] = 100;
		planned[6] = 100;
		planned[7] = 100;
		planned[8] = 100;
		planned[9] = 100;
		planned[10] = 100;
		planned[11] = 100;
		planned[12] = 100;
		planned[13] = 100;
		planned[14] = 100;
		planned[15] = 100;
		planned[16] = 100;
		planned[17] = 100;
		planned[18] = 100;
		planned[19] = 100;
		planned[20] = 100;
		planned[21] = 100;
		planned[22] = 100;
		planned[23] = 100;
		planned[24] = 100;
		planned[25] = 100;
		planned[26] = 100;
		planned[27] = 100;
		planned[28] = 100;
		planned[29] = 100;
		planned[30] = 100;
		planned[31] = 100;
		planned[32] = 100;
		planned[33] = 100;
		planned[34] = 100;
		planned[35] = 100;
		planned[36] = 100;
		planned[37] = 100;
		planned[38] = 100;
		planned[39] = 100;
		planned[40] = 100;
		planned[41] = 100;
		planned[42] = 100;
		planned[43] = 100;
		planned[44] = 100;
		planned[45] = 100;
		planned[46] = 100;
		planned[47] = 100;
		planned[48] = 100;
		planned[49] = 100;
		planned[50] = 100;
		planned[51] = 100;
		planned[52] = 100;
		planned[53] = 100;
		planned[54] = 100;
		planned[55] = 100;
		planned[56] = 100;
		planned[57] = 100;
		planned[58] = 100;
		planned[59] = 100;
		planned[60] = 100;
		planned[61] = 100;
		planned[62] = 100;
		planned[63] = 100;
		planned[64] = 100;
		planned[65] = 100;
		planned[66] = 100;
		planned[67] = 100;
		planned[68] = 100;
		planned[69] = 100;
		planned[70] = 100;
		planned[71] = 100;
		planned[72] = 100;
		planned[73] = 100;
		planned[74] = 100;
		planned[75] = 100;
		planned[76] = 100;
		planned[77] = 100;
		planned[78] = 100;
		planned[79] = 100;
		planned[80] = 100;
		planned[81] = 100;
		planned[82] = 100;
		planned[83] = 100;
		planned[84] = 100;
		planned[85] = 100;
		planned[86] = 100;
		planned[87] = 100;
		planned[88] = 100;
		planned[89] = 100;
		planned[90] = 100;
		planned[91] = 100;
		planned[92] = 100;
		planned[93] = 100;
		planned[94] = 100;
		planned[95] = 100;
		planned[96] = 100;
		planned[97] = 100;
		planned[98] = 100;
		planned[99] = 100;
		planned[100] = 100;
		planned[101] = 100;
		planned[102] = 100;
		planned[103] = 100;
		planned[104] = 100;
		planned[105] = 100;
		planned[106] = 100;
		planned[107] = 100;
		planned[108] = 100;
		planned[109] = 100;
		planned[110] = 100;
		planned[111] = 100;
		planned[112] = 100;
		planned[113] = 100;
		planned[114] = 100;
		planned[115] = 100;
		planned[116] = 100;
		planned[117] = 100;
		planned[118] = 100;
		planned[119] = 100;
		planned[120] = 100;

		uint[] memory actual = new uint[](120);
		actual[0] = 90;
		actual[1] = 90;
		actual[2] = 90;
		actual[3] = 90;
		actual[4] = 90;
		actual[5] = 90;
		actual[6] = 90;
		actual[7] = 90;
		actual[8] = 90;
		actual[9] = 90;
		actual[10] = 70;
		actual[11] = 70;
		actual[12] = 70;
		actual[13] = 70;
		actual[14] = 70;
		actual[15] = 70;
		actual[16] = 70;
		actual[17] = 70;
		actual[18] = 70;
		actual[19] = 70;
		actual[20] = 70;
		actual[21] = 70;
		actual[22] = 70;
		actual[23] = 70;
		actual[24] = 70;
		actual[25] = 70;
		actual[26] = 70;
		actual[27] = 70;
		actual[28] = 70;
		actual[29] = 70;
		actual[30] = 70;
		actual[31] = 70;
		actual[32] = 70;
		actual[33] = 70;
		actual[34] = 70;
		actual[35] = 70;
		actual[36] = 70;
		actual[37] = 70;
		actual[38] = 70;
		actual[39] = 70;
		actual[40] = 70;
		actual[41] = 70;
		actual[42] = 70;
		actual[43] = 70;
		actual[44] = 70;
		actual[45] = 70;
		actual[46] = 70;
		actual[47] = 70;
		actual[48] = 70;
		actual[49] = 70;
		actual[50] = 70;
		actual[51] = 70;
		actual[52] = 70;
		actual[53] = 70;
		actual[54] = 70;
		actual[55] = 70;
		actual[56] = 70;
		actual[57] = 70;
		actual[58] = 70;
		actual[59] = 70;
		actual[60] = 70;
		actual[61] = 70;
		actual[62] = 70;
		actual[63] = 70;
		actual[64] = 70;
		actual[65] = 70;
		actual[66] = 70;
		actual[67] = 70;
		actual[68] = 70;
		actual[69] = 70;
		actual[70] = 70;
		actual[71] = 70;
		actual[72] = 70;
		actual[73] = 70;
		actual[74] = 70;
		actual[75] = 70;
		actual[76] = 70;
		actual[77] = 70;
		actual[78] = 70;
		actual[79] = 70;
		actual[80] = 70;
		actual[81] = 70;
		actual[82] = 70;
		actual[83] = 70;
		actual[84] = 70;
		actual[85] = 70;
		actual[86] = 70;
		actual[87] = 70;
		actual[88] = 70;
		actual[89] = 70;
		actual[90] = 70;
		actual[91] = 70;
		actual[92] = 70;
		actual[93] = 70;
		actual[94] = 70;
		actual[95] = 70;
		actual[96] = 70;
		actual[97] = 70;
		actual[98] = 70;
		actual[99] = 70;
		actual[100] = 80;
		actual[101] = 80;
		actual[102] = 80;
		actual[103] = 80;
		actual[104] = 80;
		actual[105] = 80;
		actual[106] = 80;
		actual[107] = 80;
		actual[108] = 80;
		actual[109] = 80;
		actual[110] = 80;
		actual[111] = 80;
		actual[112] = 80;
		actual[113] = 80;
		actual[114] = 80;
		actual[115] = 80;
		actual[116] = 80;
		actual[117] = 80;
		actual[118] = 80;
		actual[119] = 80;
		actual[120] = 80;

		string[] memory str = new string[](2);
		str[0] = "btc miner1";
		str[1] = "btc miner2";

		address[] memory mappingAddresses = new address[](2);
		mappingAddresses[0] = 0x66b28f9dd1cf6314a8b5d691aeec6c6eaf456cbd9a;
		mappingAddresses[1] = 0x6684eb1e592de5f92059f6c7766a04f2c9fd5587ad;

		string[] memory domains = new string[](2);
		domains[0] = "mining pool1";
		domains[1] = "mining pool2";

		bool result = execute(consensus, "updateValidatorsBlockInfo(address[],uint[],uint[],string[],address[],uint,string[])", abi.encode(validators,planned,actual,str,mappingAddresses,1,domains));
		assertTrue(result);
	}

}