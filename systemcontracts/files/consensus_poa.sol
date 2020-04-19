pragma solidity 0.4.25;
pragma experimental ABIEncoderV2;

import "./library/acl.sol";
import "./library/SafeMath.sol";

/**
 * @dev POA consensus provided to business chain
 * Note that Asimov MainNet runs Satoshi+ consensus
 */

contract ConsensusPOA is ACL {
	/// delegate contract address
    address constant delegateAddr = 0x63000000000000000000000000000000000000006c;
    /// initialized or not
    bool private initialized;
    /// admins can manage both admins and validators
	address[] admins;
	/// validators in POA consensus
	address[] validators;

    mapping(address => bool) existingAdminsCheck;
    mapping(address => bool) existingValidatorsCheck;

	function init(address[] _admins, address[] _miners) public {
        require(!initialized, "it is not allowed to init more than once");

        uint adminLength = _admins.length;
        require(adminLength > 0, "invalid admin length");

        for (uint i = 0; i < adminLength; i++) {
        	existingAdminsCheck[_admins[i]] = true;
        }
        uint minerLength = _miners.length;
        for (uint j = 0; j < minerLength; j++) {
        	existingValidatorsCheck[_miners[j]] = true;
        }
        admins = _admins;
        validators = _miners;
        initialized = true;
    }

    /// batch add admins, can only be called by admins
	function batchAddAdmins(address[] newAdmins) public authAddresses(admins) {
		uint length = newAdmins.length;
		for (uint i = 0; i < length; i++) {
			if (!existingAdminsCheck[newAdmins[i]]) {
				admins.push(newAdmins[i]);
				existingAdminsCheck[newAdmins[i]] = true;
			}
		}
	}

	/// batch remove admins, can only be called by admins
	function batchRemoveAdmins(address[] oldAdmins) public authAddresses(admins) {
		uint oldAdminLength = oldAdmins.length;
		uint adminLength = admins.length;
		for (uint i = 0; i < oldAdminLength; i++) {
			for (uint j = 0; j < adminLength; j++) {
				if (oldAdmins[i] == admins[j]) {
					if (j != adminLength-1) {
						admins[j] = admins[adminLength-1];
					}
					admins.length--;
					adminLength = SafeMath.sub(adminLength, 1);
					delete existingAdminsCheck[oldAdmins[i]];
					break;
				}
			}
		}
	}

	/// batch add validators, can only be called by admins
	function batchAddValidators(address[] newValidators) public authAddresses(admins) {
		uint length = newValidators.length;
		for (uint i = 0; i < length; i++) {
			if (!existingValidatorsCheck[newValidators[i]]) {
				validators.push(newValidators[i]);
				existingValidatorsCheck[newValidators[i]] = true;
			}
		}
	}

	/// batch remove validators, can only be called by admins
	function batchRemoveValidators(address[] oldValidataors) public authAddresses(admins) {
		uint oldValidatorLength = oldValidataors.length;
		uint validatorLength = validators.length;
		for (uint i = 0; i < oldValidatorLength; i++) {
			for (uint j = 0; j < validatorLength; j++) {
				if (oldValidataors[i] == validators[j]) {
					if (j != validatorLength-1) {
						validators[j] = validators[validatorLength-1];
					}
					validators.length--;
					validatorLength = SafeMath.sub(validatorLength, 1);
					delete existingValidatorsCheck[oldValidataors[i]];
					break;
				}
			}
		}
	}

	/// transfer asset to a given address
	function transferAsset(address to, uint asset, uint amount) public authAddresses(admins) {
		require(amount > 0 && flow.balance(this, asset) >= amount, "insufficient balance");
		to.transfer(amount, asset);
	}

	/// this method is called by system; however it is only useful in Satoshi+ and it does nothing in POA
	function updateValidatorsBlockInfo(address[] validatorAddresses, uint[] plannedBlocks, uint[] actualBlocks,
        string[] validatorsStr, address[] mappingAddresses, uint round, string[] domains) public {

	}

	/// get list of validators, the parameter is only useful in Satoshi+
	function getSignupValidators(string[] validatorsStr) public view returns(address[], uint[]) {
		return (validators, new uint[](validators.length));
	}

	/// get admins and validators
	function getAdminsAndValidators() public view returns(address[], address[]) {
		return (admins, validators);
	}

}