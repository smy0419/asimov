pragma solidity 0.4.25;

import "./library/SafeMath.sol";

/**
 * @dev This system contract is designed to support the Foundation on asimov website
 */

contract GenesisOrganization {
    /// delegate contract address
    address constant delegateAddr = 0x630000000000000000000000000000000000000064;
    /// initialized or not
    bool private initialized;
    /// proposal types
    enum ProposalType {ELECT, IMPEACH, EXPENSES}
    /// proposal status
    enum ProposalStatus {ONGOING, APPROVED, REJECTED}

    /// member of genesis organization
    struct Citizen {
        mapping(uint => bool) votedProposals;
        bool existed;
    }

    /// proposal data structure
    struct Proposal {
        address proposer;
        address[] voters;
        address[] approvers;
        address[] rejecters;
        uint voterLength;
        uint approverLength;
        uint rejecterLength;
        uint percent;
        uint startHeight;
        uint endHeight;
        ProposalType proposalType;
        ProposalStatus status;
        bool existed;

        /// only useful to EXPENSES
        address targetAddress;
        uint asset;
        uint amount;
    }

    /// address => Citizen
    mapping(address => Citizen) citizenMap;
    /// proposal id => Proposal
    mapping(uint => Proposal) proposalMap;
    /// map for quick check
    mapping(address => bool) existingCitizensCheck;

    /// all members
    address[] citizens;
    /// proposal index, auto incremental
    uint private proposalIndex;
    /// minimal members required to run the organization
    uint private MINIMAL_CITIZEN_COUNT;

    event ReceiveEvent(address donator, uint assetType, uint amount);
    event TransferAssetEvent(address receivingAddress, ProposalType proposalType, uint assetType, uint amount);
    event StartProposalEvent(uint proposalId, ProposalType proposalType, address proposer, uint endTime);
    event VoteProposalEvent(uint proposalId, address[] voters);
    event VoteEvent(uint proposalId, address voter, bool decision);
    event ProposalStatusChangeEvent(uint proposalId, ProposalStatus status);
    event ElectCitizenEvent(address newCitizen);
    event ImpeachCitizenEvent(address oldCitizen);

    /// payable fallback function
    function() public payable {
        emit ReceiveEvent(msg.sender, msg.assettype, msg.value);
    }

    function init(address[] genesisCitizens) public {
        require(!initialized, "it is not allowed to init more than once");
        uint length = genesisCitizens.length;
        require(length == 7, "require 7 members");

        for (uint i = 0; i < length; i++) {
            Citizen storage citizen = citizenMap[genesisCitizens[i]];
            citizen.existed = true;
            existingCitizensCheck[genesisCitizens[i]] = true;
        }
        citizens = genesisCitizens;

        MINIMAL_CITIZEN_COUNT = 3;
        initialized = true;
    }

    /**
     * @dev make a proposal, can only be called by a member
     *
     * @param proposalType proposal type
     * @param targetAddress target address to receive asset (only useful for an expense proposal)
     * @param amount asset amount (only useful for an expense proposal)
     * @param asset asset type (only useful for an expense proposal)
     * @return proposal id
     */
    function startProposal(uint proposalType, address targetAddress, uint amount, uint asset) public returns(uint) {
        require(proposalType >= 0 && proposalType <= 2, "invalid proposal type");
        require(existingCitizensCheck[msg.sender], "not authorized");

        proposalIndex = SafeMath.add(proposalIndex, 1);
        Proposal storage prop = proposalMap[proposalIndex];
        prop.proposer = msg.sender;
        prop.voters = citizens;
        prop.voterLength = citizens.length;
        prop.targetAddress = targetAddress;
        prop.startHeight = block.number;
        /// 7 days, calculated in blocks
        prop.endHeight = SafeMath.add(block.number, SafeMath.mul(SafeMath.mul(7, 24), 720));
        prop.status = ProposalStatus.ONGOING;
        prop.existed = true;

        // elect new member
        if (uint(ProposalType.ELECT) == proposalType) {
            require(!existingCitizensCheck[targetAddress], "already a member");
            prop.percent = 70;
            prop.proposalType = ProposalType.ELECT;
        }
        // impeach a member
        else if (uint(ProposalType.IMPEACH) == proposalType) {
            require(existingCitizensCheck[targetAddress], "not a member");
            require(citizens.length >= MINIMAL_CITIZEN_COUNT, "cannot do that, not enough members left");
            prop.percent = 70;
            prop.proposalType = ProposalType.IMPEACH;
        }
        // expense
        else if (uint(ProposalType.EXPENSES) == proposalType) {
            require(amount > 0, "insufficient balance");
            require(asset & 0x10000000000000000 == 0 ? flow.balance(this, asset) >= amount : flow.balance(this, asset) > 0, "insufficient balance");
            prop.asset = asset;
            prop.amount = amount;
            prop.percent = 50;
            prop.proposalType = ProposalType.EXPENSES;
        }

        emit StartProposalEvent(proposalIndex, prop.proposalType, msg.sender, SafeMath.add(block.timestamp, SafeMath.mul(SafeMath.mul(7, 24), 3600)));
        emit VoteProposalEvent(proposalIndex, prop.voters);

        return proposalIndex;
    }

    /**
     * @dev vote on a proposal, can only be called by a member
     *
     * @param proposalId proposal id
     * @param decision yes or no
     */
    function vote(uint proposalId, bool decision) public {
        require(existingCitizensCheck[msg.sender], "not authorized");

        Proposal storage prop = proposalMap[proposalId];
        require(prop.existed, "proposal does not exist");
        require(prop.endHeight >= block.number, "vote expired");
        require(prop.status == ProposalStatus.ONGOING, "invalid proposal status");

        Citizen storage citizen = citizenMap[msg.sender];
        require(!citizen.votedProposals[proposalId], "already voted");

        emit VoteEvent(proposalId, msg.sender, decision);

        if (decision) {
            prop.approvers.push(msg.sender);
            prop.approverLength = SafeMath.add(prop.approverLength, 1);
            if (SafeMath.div(SafeMath.mul(prop.approverLength, 100), prop.voterLength) >= prop.percent) {
                prop.status = ProposalStatus.APPROVED;
                emit ProposalStatusChangeEvent(proposalId, prop.status);

                executeProposal(prop.proposalType, prop.targetAddress, prop.asset, prop.amount);
            }
        } else {
            prop.rejecters.push(msg.sender);
            prop.rejecterLength = SafeMath.add(prop.rejecterLength, 1);
            if (SafeMath.div(SafeMath.mul(prop.rejecterLength, 100), prop.voterLength) > SafeMath.sub(100, prop.percent)) {
                prop.status = ProposalStatus.REJECTED;
                emit ProposalStatusChangeEvent(proposalId, prop.status);
            }
        }
        citizen.votedProposals[proposalId] = true;
    }
    
    /**
     * execute a proposal
     *
     * @param proposalType proposal type
     * @param targetAddress target address to receive asset (only useful for an expense proposal)
     * @param asset asset type (only useful for an expense proposal)
     * @param amount asset amount (only useful for an expense proposal)
     */
    function executeProposal(ProposalType proposalType, address targetAddress, uint asset, uint amount) internal {
        // add new member
        if (ProposalType.ELECT == proposalType) {
            if (!existingCitizensCheck[targetAddress]) {
                Citizen storage citizen = citizenMap[targetAddress];
                citizen.existed = true;
                citizens.push(targetAddress);
                existingCitizensCheck[targetAddress] = true;
                emit ElectCitizenEvent(targetAddress);
            }
        }
        // remove a member
        else if (ProposalType.IMPEACH == proposalType) {
            require(existingCitizensCheck[targetAddress], "target address may had impeached");
            require(citizens.length >= MINIMAL_CITIZEN_COUNT, "cannot do that, not enough members left");
            uint length = citizens.length;
            for (uint i = 0; i < length; i++) {
                if (targetAddress == citizens[i]) {
                    if (i != length-1) {
                        citizens[i] = citizens[length-1];
                    }
                    citizens.length--;
                    emit ImpeachCitizenEvent(targetAddress);
                    break;
                }
            }
            delete existingCitizensCheck[targetAddress];
        }
        // make the expense
        else if (ProposalType.EXPENSES == proposalType) {
            targetAddress.transfer(amount, asset);
            emit TransferAssetEvent(targetAddress, ProposalType.EXPENSES, asset, amount);
        }
    }

    /**
     * check if caller is one of citizens
     *
     * @param caller caller address
     * @return if the caller is one of citizens
     */
    function existed(address caller) public returns (bool) {
        return existingCitizensCheck[caller];
    }


}