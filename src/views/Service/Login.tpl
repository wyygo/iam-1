{{template "Common/HtmlHeader.tpl" .}}

<style>
body {
  margin: 0 auto;
  position: relative;
}

#ids-login-ctn {
  width: 320px;
  position: absolute;
  left: 50%;
  top: 20%;
  margin-left: -160px;
  color: #555;
}

#ids-login-form {
  background-color: #f7f7f7;
  border-radius: 4px;
  padding: 40px;
  box-shadow: 0px 2px 2px 0px #999;
}

.ids-login-msg01 {
  font-size: 18px;
  padding: 10px 0;
  text-align: center;
}

.ids-user-ico-default {
  width: 80px;
  height: 80px;
  position: relative;
  left: 50%;
  margin: 0 0 30px -40px;
}

#ids-login-form .ilf-group {
  padding: 0 0 10px 0; 
}

#ids-login-form .ilf-input {
  display: block;
  width: 100%;
  height: 40px;
  padding: 5px 10px;
  font-size: 14px;
  line-height: 1.42857143;
  color: #555;
  background-color: #fff;
  background-image: none;
  border: 1px solid #ccc;
  border-radius: 2px;
  box-shadow: inset 0 1px 1px rgba(0, 0, 0, .075);
  transition: border-color ease-in-out .15s, box-shadow ease-in-out .15s;
}

#ids-login-form .ilf-input:focus {
  border-color: #66afe9;
  outline: 0;
  box-shadow: inset 0 1px 1px rgba(0,0,0,.075), 0 0 8px rgba(102, 175, 233, .6);
}

#ids-login-form .ilf-btn {
  width: 100%;
  height: 40px;
  display: inline-block;
  padding: 5px 10px;
  margin-bottom: 0;
  font-size: 14px;
  font-weight: normal;
  line-height: 1.42857143;
  text-align: center;
  white-space: nowrap;
  vertical-align: middle;
  cursor: pointer;
  -webkit-user-select: none;
     -moz-user-select: none;
      -ms-user-select: none;
          user-select: none;
  background-image: none;
  border: 1px solid transparent;
  border-radius: 3px;
  color: #fff;
  background-color: #428bca;
  border-color: #357ebd;
}

#ids-login-form .ilf-btn:hover {
  color: #fff;
  background-color: #3276b1;
  border-color: #285e8e;
}

#ids-login-ctn .il-footer {
  text-align: center;
  margin: 10px 0;
  font-size: 14px;
}
#ids-login-ctn .il-footer img {
  width: 16px;
  height: 16px;
}
</style>

<div id="ids-login-ctn">

    <div class="ids-login-msg01">{{T . "Sign in with your Account"}}</div>

    <form id="ids-login-form" class="" action="#">

      <img class="ids-user-ico-default"  src="/ids/static/img/user-default.png">

      <div class="alert alert-info hide"></div>

      <div class="ilf-group">
        <input type="text" class="ilf-input" name="userid" placeholder="{{T . "Username"}}">
      </div>

      <div class="ilf-group">
        <input type="password" class="ilf-input" name="passwd" placeholder="{{T . "Password"}}">
      </div>

      <div class="">
        <button type="submit" class="ilf-btn">{{T . "Sign in"}}</button>
      </div>
    </form>

    <div class="il-footer">
      <img src="/ids/static/img/ids-s1-32.png"> 
      <a href="http://www.lesscompute.com" target="_blank">less Identity Server</a>
    </div>
</div>

<script>


var ids_eh = $("#ids-login-ctn").height();
$("#ids-login-ctn").css({
    "top": "45%",
    "margin-top": - (ids_eh / 2) + "px" 
});

$("input[name=userid]").focus();

$("#gbfg5g").submit(function(event) {

    event.preventDefault();

    var req = {
      data: {
        "userid": $("input[name=userid]").val(),
        "passwd": $("input[name=passwd]").val(),
      }
    }

    //console.log(JSON.stringify(req));

    $.ajax({
        type    : "POST",
        url     : "/ids/service/auth",
        data    : JSON.stringify(req),
        timeout : 3000,
        contentType: "application/json; charset=utf-8",
        success : function(rsp) {

            var rsj = JSON.parse(rsp);
            //console.log(rsp);

            if (rsj.status == 200) {
                lessCookie.Set("access_token_lessfly", rsj.access_token, 7200);
                $('#body-content').load('/lessfly/index/well');
                //saComLoader('index/index');
            } else {
                alert(rsj.message);
            }
        },
        error: function(xhr, textStatus, error) {
            alert('{{T . "Internal Server Error"}}');
        }
    });
});


</script>
{{template "Common/HtmlFooter.tpl" .}}