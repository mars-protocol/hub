---- MODULE envoy ----
EXTENDS Apalache, TLC, Integers, Sequences, SequencesExt, Variants

(*
This TLA+ specification describes the behavior of the envoy module implementation as part of the Mars Protocol.
The module has a state, typically called Keeper in Cosmos-SDK convention.
There are three ways to make changes in the module state - executing three types of transactions defined by the module.

For more details of the system behaviors and implementation, read the System Behavior section in [Informal System audit report from Q1 2023](https://github.com/informalsystems/audit-mars-protocol/blob/main/2023/Q1/report/report.pdf).

The specification models the Keeper using state variables.
It also models the executed transaction and its success using an addition variable `action`.

The init and next operators as defined as `Init` and `Next`.
`Next` operator is a disjunction of three specific next operators corresponding to the three transaction types.
There are also few invariants at the end of the specification.
*)

(*

@typeAlias: remoteMsg = Str;
@typeAlias: connectionId = Int;
@typeAlias: channelId = Int;
@typeAlias: denom = Str;
@typeAlias: accountId = Str;
@typeAlias: portId = Str;
@typeAlias: balance = $denom -> Int;

@typeAlias: bankKeeper = $accountId -> $balance;
@typeAlias: channelKeeper = $channelId -> {connection_id: $connectionId, port: $portId};
@typeAlias: icaControllerKeeper = $connectionId -> $channelId;

@typeAlias: ibcPacket =
    SendFunds({remote: $accountId, amount: $balance}) |
    SendMessages({remote: $accountId, messages: Seq($remoteMsg)});

@typeAlias: msg =
    RegisterAccount({connection_id: Int}) |
    SendFunds({authority: Str, channel_id: Int, amount: $balance}) |
    SendMessages({authority: Str, connection_id: Int, messages: Seq($remoteMsg)}) |
    Genesis(Int);

@typeAlias: trace = {
    bank_keeper: $bankKeeper,
    channel_keeper: $channelKeeper,
    ica_controller_keeper: $icaControllerKeeper,
    ibc_packets: Seq($ibcPacket),
    authority: $accountId,
    action: {msg: $msg, success: Bool}
};

*)

GOV_ACCOUNT == "gov"
FEE_POOL_ACCOUNT == "fee_pool"
IBC_ESCROW_ACCOUNT == "ibc_escrow"
ENVOY_ACCOUNT == "envoy"

ACCOUNTS == {GOV_ACCOUNT, FEE_POOL_ACCOUNT, IBC_ESCROW_ACCOUNT, ENVOY_ACCOUNT, "Alice", "Bob"}

DENOMS == {"umars", "uosmo"}

\* connection-0
CONNECTION_ID == 0

REMOTE_MSGS == {"bank/send", "cw/update"}

IBC_TRANSFER_PORT == "transfer"
ICA_CONTROLLER_PORT == "ica-controller"


VARIABLES
    \* @type: $bankKeeper;
    bank_keeper,

    \* @type: $channelKeeper;
    channel_keeper,

    \* @type: $icaControllerKeeper;
    ica_controller_keeper,

    \* @type: Seq($ibcPacket);
    ibc_packets,

    \* @type: $accountId;
    authority,

    \* @type: {msg: $msg, success: Bool};
    action


\* The genesis creation logic is implemented in [app.go](https://github.com/mars-protocol/hub/blob/c7795c2488447f4d2683883277f75c55b3505f03/app/app.go)
Init ==
    \E _channel_id \in Nat:
    \E _bank_keeper \in [ACCOUNTS -> [DENOMS -> 0..10]]:
        /\ bank_keeper = _bank_keeper
        /\ channel_keeper = SetAsFun({<<_channel_id, [connection_id |-> CONNECTION_ID, port |-> IBC_TRANSFER_PORT]>>})
        /\ ica_controller_keeper = SetAsFun({})
        /\ ibc_packets = <<>>
        \* Envoy authority is [set to gov module account](https://github.com/mars-protocol/hub/blob/c7795c2488447f4d2683883277f75c55b3505f03/app/app.go#L337)
        /\ authority = GOV_ACCOUNT
        /\ action = [msg |-> Variant("Genesis", 0), success |-> TRUE]


\* [RegisterAccount](https://github.com/mars-protocol/hub/blob/c7795c2488447f4d2683883277f75c55b3505f03/x/envoy/keeper/msg_server.go#L88)
RegisterAccountNext ==
    \E _connection_id \in Nat:
    \E _new_channel_id \in Nat:
        LET
        _msg == Variant("RegisterAccount", [connection_id |-> _connection_id])
        _is_success ==
            \* the connection ID must be active
            /\ _connection_id \in {CONNECTION_ID}
            \* there should not be an existing ICA
            /\ _connection_id \notin DOMAIN ica_controller_keeper
            \* the new channel ID must be unused
            /\ _new_channel_id \notin DOMAIN channel_keeper
        IN
            IF _is_success THEN
                /\ channel_keeper' = channel_keeper @@ (_new_channel_id :> [connection_id |-> _connection_id, port |-> ICA_CONTROLLER_PORT])
                /\ ica_controller_keeper' = ica_controller_keeper @@ (_connection_id :> _new_channel_id)
                /\ action' = [msg |-> _msg, success |-> TRUE]
                /\ UNCHANGED <<bank_keeper, ibc_packets, authority>>
            ELSE
                /\ action' = [msg |-> _msg, success |-> FALSE]
                /\ UNCHANGED <<bank_keeper, channel_keeper, ica_controller_keeper, ibc_packets, authority>>


\* @type: ($balance) => Bool;
IsEmpty(_amount) ==
    LET
    \* @type: (Bool, Bool) => Bool;
    NotLambda(_x, _y) == _x /\ _y
    IN ApaFoldSet(NotLambda, TRUE, {_amount[_x] <= 0: _x \in DOMAIN _amount})


\* @type: ($balance) => Bool;
AllPositiveAmount(_amount) ==
    LET
    \* @type: (Bool, Bool) => Bool;
    NotLambda(_x, _y) == _x /\ _y
    IN ApaFoldSet(NotLambda, TRUE, {_amount[_x] >= 0: _x \in DOMAIN _amount})


\* @type: (Int, Int) => Int;
Min(a, b) ==
    IF a < b THEN a ELSE b


\* @type: (Int, Int) => Int;
Max(a, b) ==
    IF a > b THEN a ELSE b


\* [SendFunds](https://github.com/mars-protocol/hub/blob/c7795c2488447f4d2683883277f75c55b3505f03/x/envoy/keeper/msg_server.go#L118)
SendFundsNext ==
    \E _account \in DOMAIN bank_keeper:
    \E _channel_id \in Nat:
    \E _amount \in [DENOMS -> Nat]:
        LET
        _msg == Variant("SendFunds", [authority |-> _account, channel_id |-> _channel_id, amount |-> _amount])
        _old_envoy_balance == bank_keeper[ENVOY_ACCOUNT]
        _new_envoy_balance == [_d \in DENOMS |-> Max(_old_envoy_balance[_d] - _amount[_d], 0)]
        _envoy_short_fall == [_d \in DENOMS |-> Min(_old_envoy_balance[_d] - _amount[_d], 0)]
        _old_fee_pool_balance == bank_keeper[FEE_POOL_ACCOUNT]
        _new_fee_pool_balance == [_d \in DENOMS |-> Max(_old_fee_pool_balance[_d] + _envoy_short_fall[_d], 0)]
        _fee_pool_short_fall == [_d \in DENOMS |-> Min(_old_fee_pool_balance[_d] + _envoy_short_fall[_d], 0)]
        _old_ibc_escrow_balance == bank_keeper[IBC_ESCROW_ACCOUNT]
        _new_ibc_escrow_balance == [_d \in DENOMS |-> _old_ibc_escrow_balance[_d] + _amount[_d]]
        _is_success ==
            \* the sent funds must be non-empty
            /\ ~IsEmpty(_amount)
            \* the accounts must have enough balance for the transfer
            /\ AllPositiveAmount(_fee_pool_short_fall)
            \* the submitter must be the authority
            /\ _account = authority
            \* the ICAccount must exist
            /\ channel_keeper[_channel_id].connection_id \in DOMAIN ica_controller_keeper
            \* the channel_id must be active
            /\ _channel_id \in DOMAIN channel_keeper
            \* the channel id must be for ibc-transfer
            /\ channel_keeper[_channel_id].port = IBC_TRANSFER_PORT
        IN
        IF _is_success THEN
            /\ bank_keeper' = [bank_keeper EXCEPT
                ![ENVOY_ACCOUNT] = _new_envoy_balance,
                ![FEE_POOL_ACCOUNT] = _new_fee_pool_balance,
                ![IBC_ESCROW_ACCOUNT] = _new_ibc_escrow_balance
                ]
            /\ action' = [msg |-> _msg, success |-> TRUE]
            /\ ibc_packets' = Append(ibc_packets, Variant("SendFunds", [remote |-> ENVOY_ACCOUNT, amount |-> _amount]))
            /\ UNCHANGED <<channel_keeper, ica_controller_keeper, authority>>
        ELSE
            action' = [msg |-> _msg, success |-> FALSE]
            /\ UNCHANGED <<bank_keeper, channel_keeper, ica_controller_keeper, ibc_packets, authority>>


\* [SendMessages](https://github.com/mars-protocol/hub/blob/c7795c2488447f4d2683883277f75c55b3505f03/x/envoy/keeper/msg_server.go#L210)
SendMessagesNext ==
    \E _account \in DOMAIN bank_keeper:
    \E _connection_id \in Nat:
    \E _remote_msgs_set \in SUBSET REMOTE_MSGS:
        LET
        _remote_msgs == SetToSeq(_remote_msgs_set)
        _msg == Variant("SendMessages", [authority |-> _account, connection_id |-> _connection_id, messages |-> _remote_msgs])
        _is_success ==
            \* the submitter must be the authority
            /\ _account = authority
            \* the ICAccount must exist
            /\ _connection_id \in DOMAIN ica_controller_keeper
            \* the channel id must be for ICA
            /\ channel_keeper[ica_controller_keeper[_connection_id]].port = ICA_CONTROLLER_PORT
        IN
        IF _is_success THEN
            /\ action' = [msg |-> _msg, success |-> TRUE]
            /\ ibc_packets' = Append(ibc_packets, Variant("SendMessages", [remote |-> ENVOY_ACCOUNT, messages |-> _remote_msgs]))
            /\ UNCHANGED <<bank_keeper, channel_keeper, ica_controller_keeper, authority>>
        ELSE
            action' = [msg |-> _msg, success |-> FALSE]
            /\ UNCHANGED <<bank_keeper, channel_keeper, ica_controller_keeper, ibc_packets, authority>>

\* Specifies the behaviors implemented in [/x/envoy/keeper/msg_server.go](https://github.com/mars-protocol/hub/blob/c7795c2488447f4d2683883277f75c55b3505f03/x/envoy/keeper/msg_server.go)
Next ==
    \/ RegisterAccountNext
    \/ SendFundsNext
    \/ SendMessagesNext

(*

A trace is a sequence of valid states starting with the _init_ operator and then applying the _next_ operator on two consequent operators.

```
init(init[0])
next(init[0], init[1])
next(init[1], init[2])
...
next(init[i], init[i+1])
```

A property is a boolean operator that the current trace or the current state and asserts a property on them.

Here we list some properties as `invariant-xx` and `non-invariant-xx`.

`invariant-xx` are properties that are asserted true by all possible traces.

`non-invariant-xx` are properties that can be asserted false by atleast one trace.

`non-invariant` properties are used to generate example traces, that can be used to test the original implementation.
They are usually negation of another property, `example-xx`.
So a counter example of a `non-invariant` satisfies the negation of the negation of `example-xx` or simply, it satisfies `example-xx`.

There are some other operators, mentioned as view operator, that are required for multiple example generation.
Each example is unique when projected using the view operator.

*)


\* example-AS
\* Asserts the trace has minimum 5 states and all the actions succeeded.
\* @type: Seq($trace) => Bool;
AllSuccess(_trace) ==
    /\ Len(_trace) >= 5
    /\ \A _i \in DOMAIN _trace: _trace[_i].action.success


\* non-invariant-AS
\* Negation of example-AS; AllSuccess
\* This will produce counter examples whose all actions succeeded.
\* @type: Seq($trace) => Bool;
ExAllSuccess(_trace) == ~AllSuccess(_trace)

\* view-AT
\* Projects a state to action message type.
\* () => Str;
ActionType == VariantTag(action.msg)

\* invariant-CS
\* Asserts the token supply on the chain never changes.
\* @type: Seq($trace) => Bool;
ConstantSupply(_trace) ==
    \A _i, _j \in DOMAIN _trace:
    \A _d \in DENOMS:
        LET
        \* @type: (Int, <<$accountId, Int>>) => Int;
        NotLambda(_x, _y) == _x + _y[2]
        \* @type: ($bankKeeper) => Int;
        DenomSupply(_bank_keeper) == ApaFoldSet(NotLambda, 0, {<<_acc, _bank_keeper[_acc][_d]>>: _acc \in DOMAIN _bank_keeper})
        IN
        DenomSupply(_trace[_i].bank_keeper) = DenomSupply(_trace[_j].bank_keeper)


\* invariant-PB
\* Asserts balances must always be non-negative.
\* @type: Seq($trace) => Bool;
PositiveBalance(_trace) ==
    \A _i \in DOMAIN _trace:
        LET
        _state == _trace[_i]
        \* @type: (Bool, $balance) => Bool;
        NotLambda(_x, _y) == _x /\ AllPositiveAmount(_y)
        IN
        ApaFoldSet(NotLambda, TRUE, {_state.bank_keeper[_acc]: _acc \in DOMAIN _state.bank_keeper})


\* invariant-VA
\* Asserts the successful SendFunds and SendMessages have the correct authority.
\* @type: Seq($trace) => Bool;
ValidAuthority(_trace) ==
    \A _i \in DOMAIN _trace:
        LET _state == _trace[_i] IN
        CASE VariantTag(_state.action.msg)  = "SendFunds" ->
            LET _msg == VariantGetUnsafe("SendFunds", _state.action.msg) IN
            _msg.authority /= GOV_ACCOUNT
            =>
            ~_state.action.success
        [] VariantTag(_state.action.msg)  = "SendMessages" ->
            LET _msg == VariantGetUnsafe("SendMessages", _state.action.msg) IN
            _msg.authority /= GOV_ACCOUNT
            =>
            ~_state.action.success
        [] OTHER ->
            TRUE


\* invariant-IE
\* Asserts the successful SendFunds and SendMessages are submitted to an existing ICAccount.
\* @type: Seq($trace) => Bool;
ICAExists(_trace) ==
    \A _i \in DOMAIN _trace:
        LET _state == _trace[_i] IN
        CASE VariantTag(_state.action.msg)  = "SendFunds" ->
            _state.action.success
            =>
            LET _msg == VariantGetUnsafe("SendFunds", _state.action.msg) IN
            /\ _msg.channel_id \in DOMAIN _state.channel_keeper
            /\ _state.channel_keeper[_msg.channel_id].connection_id \in DOMAIN _state.ica_controller_keeper
        [] VariantTag(_state.action.msg)  = "SendMessages" ->
            _state.action.success
            =>
            LET _msg == VariantGetUnsafe("SendMessages", _state.action.msg) IN
            _msg.connection_id \in DOMAIN _state.ica_controller_keeper
        [] OTHER ->
            TRUE


\* invariant-PL
\* Asserts the ibc packets from successful SendFunds and SendMessages are submitted in ibc queue.
\* @type: Seq($trace) => Bool;
NoPacketLoss(_trace) ==
    \A _i \in DOMAIN _trace:
        LET _state == _trace[_i] IN
        CASE VariantTag(_state.action.msg)  = "SendFunds" ->
            _state.action.success
            =>
            LET
            _msg == VariantGetUnsafe("SendFunds", _state.action.msg)
            _ibc_msg == VariantGetUnsafe("SendFunds", Last(_state.ibc_packets))
            IN
            _msg.amount = _ibc_msg.amount
        [] VariantTag(_state.action.msg)  = "SendMessages" ->
            _state.action.success
            =>
            LET
            _msg == VariantGetUnsafe("SendMessages", _state.action.msg)
            _ibc_msg == VariantGetUnsafe("SendMessages", Last(_state.ibc_packets))
            IN
            _msg.messages = _ibc_msg.messages
        [] OTHER ->
            TRUE


\* all invariants
\* @type: Seq($trace) => Bool;
InvAll(_trace) ==
    /\ ConstantSupply(_trace)
    /\ PositiveBalance(_trace)
    /\ ValidAuthority(_trace)
    /\ ICAExists(_trace)
    /\ NoPacketLoss(_trace)

====
