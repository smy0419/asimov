pragma solidity 0.4.25;
pragma experimental ABIEncoderV2;

import "./library/SafeMath.sol";

/**
 * @dev validator committee interface
 */
interface Committee {
    function startNewRound(address[] validatorAddresses, uint[] plannedBlocks, uint[] actualBlocks) external;
}

/**
 * @dev Asimov MainNet runs Satoshi+ consensus
 */
contract ConsensusSatoshiPlus {
    /// delegate contract address
    address constant delegateAddr = 0x63000000000000000000000000000000000000006a;
    /// initialized or not
    bool private initialized;

    /// validator data structure
    struct Validator {
        /// number of blocks plan to produce
        uint plannedBlocks;
        /// number of blocks actually produced
        uint actualBlocks;
        /// efficiency = actualBlocks / plannedBlocks
        uint efficiency;
        /// which round
        uint round;
        /// existed or not
        bool existed;
    }

    /// address mapping between bitcoin and asimov
    struct MappedValidator {
        /// asimov address
        address validator;
        /// the round when asimov detects the special mapping transaction is sent on bitcoin
        uint round;
        /// existed or not
        bool existed;
    }

    /// address => Validator
    mapping(address => Validator) validators;
    /// btc address => MappedValidator
    mapping(string => MappedValidator) mappedValidators;

    mapping(address => bool) signupValidatorsCheck;
    mapping(address => bool) mappedValidatorsCheck;

    /// validator committee contract
    Committee committee;

    event MapValidatorsEvent(string[] btcAddresses, address[] asimovAddresses, string[] domains);

    function init() public {
        require(!initialized, "it is not allowed to init more than once");

        committee = Committee(0x63000000000000000000000000000000000000006b);
        initialized = true;
    }

    /**
     * @dev validators to signup for Satoshi+ consensus
     */
    function signup() public {
        require(mappedValidatorsCheck[msg.sender], "no mapped btc address");

        Validator storage validator = validators[msg.sender];
        validator.round = block.round;
        validator.existed = true;

        signupValidatorsCheck[msg.sender] = true;
    }

    /**
     * @dev get signup validator addresses before starting new round
     *
     * @param btcAddresses string format of btc addresses
     */
    function getSignupValidators(string[] btcAddresses) public view returns(address[], uint[]) {
        uint length = btcAddresses.length;
        address[] memory mappedAddresses = new address[](length);
        uint[] memory rounds = new uint[](length);

        for (uint i = 0; i < length; i++) {
            MappedValidator storage mValidator = mappedValidators[btcAddresses[i]];
            if (!mValidator.existed) {
                mappedAddresses[i] = 0x0;
                rounds[i] = 0xffffffffffffffff;
            } else {
                mappedAddresses[i] = mValidator.validator;
                rounds[i] = mValidator.round;
                if (signupValidatorsCheck[mValidator.validator]) {
                    Validator storage validator = validators[mValidator.validator];
                    /// signup round vs. mapped round
                    if (validator.round > mValidator.round) {
                        rounds[i] = validator.round;
                    }
                }
            }
        }
        return (mappedAddresses, rounds);
    }

    /**
     * @dev update validators's blocks information and new mapping relationships
     *
     * @param validatorAddresses validator addresses
     * @param plannedBlocks blocks planned to produce by each validator
     * @param actualBlocks blocks actually produced by each validator
     * @param btcAddresses string format of btc addresses
     * @param asimovAddresses asimov addresses
     * @param round which round
     * @param domains names of validators
     */
    function updateValidatorsBlockInfo(address[] validatorAddresses, uint[] plannedBlocks, uint[] actualBlocks,
        string[] btcAddresses, address[] asimovAddresses, uint round, string[] domains)
        public
    {
        require(msg.sender == 0x660000000000000000000000000000000000000000, "invalid caller");
        uint length = validatorAddresses.length;
        require(length == plannedBlocks.length && length == actualBlocks.length, "invalid paramenter length");
        uint strLength = btcAddresses.length;
        require(strLength == asimovAddresses.length && strLength == domains.length, "invalid parameter length");

        /// update validators' blocks information
        for (uint i = 0; i < length; i++) {
            Validator storage validator = validators[validatorAddresses[i]];
            validator.plannedBlocks = SafeMath.add(validator.plannedBlocks, plannedBlocks[i]);
            validator.actualBlocks = SafeMath.add(validator.actualBlocks, actualBlocks[i]);
            validator.efficiency = SafeMath.div(SafeMath.mul(validator.actualBlocks, 100), validator.plannedBlocks);
            validator.existed = true;
        }

        /// record mapping relationship of validators
        for (uint j = 0; j < strLength; j++) {
            MappedValidator storage mValidator = mappedValidators[btcAddresses[j]];
            mValidator.validator = asimovAddresses[j];
            mValidator.round = round;
            mValidator.existed = true;
            mappedValidatorsCheck[asimovAddresses[j]] = true;
        }
        emit MapValidatorsEvent(btcAddresses, asimovAddresses, domains);

        /// call startNewRound method of validator committee contract
        committee.startNewRound(validatorAddresses, plannedBlocks, actualBlocks);
    }

    /**
     * @dev get validator block information
     *
     * @param validatorAddress validator address
     * @return {planned blocks, actual blocks, efficiency}
     */
    function getValidatorBlockInfo(address validatorAddress) public view returns(uint, uint, uint) {
        Validator storage validator = validators[validatorAddress];
        if (!validator.existed) {
            return (0, 0, 0);
        }

        return (validator.plannedBlocks, validator.actualBlocks, validator.efficiency);
    }

}