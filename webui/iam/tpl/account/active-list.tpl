<div class="iam-div-light">
  <table class="table table-hover">
    <thead>
      <tr>
        <th>ID</th>
        <th>Type</th>
        <th>Amount</th>
        <th>Pre Paying</th>
        <th>Payout</th>
		<th>Product Limits</th>
		<th>Product Max</th>
		<th>Product Inpay</th>
        <th>Actived</th>
        <th></th>
      </tr>
    </thead>
    <tbody id="iam-acc-activelist"></tbody>
  </table>
</div>

<script id="iam-acc-activelist-tpl" type="text/html">
{[~it.items :v]}
<tr>
  <td class="iam-monofont">
    {[=v.id.substr(8)]}
  </td>
  <td>
    {[~it._ecoin_types :sv]}
    {[ if (v.type == sv.value) { ]}{[=sv.name]}{[ } ]}
    {[~]}
  </td>
  <td>{[=v.amount]}</td>
  <td>{[=v.prepay]}</td>
  <td>{[=v.payout]}</td>
  <td>{[=v._exp_product_limits]}</td>
  <td>{[=v.exp_product_max]}</td>
  <td>
    {[~v.exp_product_inpay :pv]}
	<div>{[=pv]}</div>
    {[~]}
  </td>
  <td>{[=l4i.MetaTimeParseFormat(v.created, "Y-m-d H:i")]}</td>
  <td align="right">
    <!--
	<button class="pure-button button-xsmall"
      onclick="iamAccessKey.Set('{[=v.access_key]}')">
      Setting
    </button>
	-->
  </td>
</tr>
{[~]}
</script>

<script type="text/html" id="iam-acc-activelist-optools">
<li class="iam-btn iam-btn-primary">
  <a href="#" onclick="iamAccessKey.Set()">
     ECoin Recharge
  </a>
</li>
</script>
