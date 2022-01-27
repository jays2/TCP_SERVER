<template>
    <!DOCTYPE html>
<html lang="en">
<title>Server stats</title>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="stylesheet" href="https://www.w3schools.com/w3css/4/w3.css">
<link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Poppins">

<body>

<!-- Sidebar/menu -->
<nav class="w3-sidebar w3-red w3-collapse w3-top w3-large w3-padding" style="z-index:3;width:300px;font-weight:bold;" id="mySidebar"><br>
  <div class="w3-container">
    <h3 class="w3-padding-64"><b>@PlainLogs.com</b></h3>
  </div>
  <div class="w3-bar-block">
    <a href="https://github.com/jays2/general_v1.git" onclick="w3_close()" class="w3-bar-item w3-button w3-hover-white">@Github</a> 
    <a href="https://www.linkedin.com/in/enid-sierra-vargas/" onclick="w3_close()" class="w3-bar-item w3-button w3-hover-white">@Copyright 2022</a> 
  </div>
</nav>

<!-- Top menu on small screens -->
<header class="w3-container w3-top w3-hide-large w3-red w3-xlarge w3-padding">
  <a href="javascript:void(0)" class="w3-button w3-red w3-margin-right" onclick="w3_open()">☰</a>
  <span>Company Name</span>
</header>

<!-- Overlay effect when opening sidebar on small screens -->
<div class="w3-overlay w3-hide-large" onclick="w3_close()" style="cursor:pointer" title="close side menu" id="myOverlay"></div>

<!-- !PAGE CONTENT! -->
<div class="w3-main" style="margin-left:340px;margin-right:40px">

  
  <!-- Channels -->
  <div class="w3-container" id="designers" style="margin-top:75px">
    <h1 class="w3-xxxlarge w3-text-red"><b>CHANNELS</b></h1>
    <hr style="width:50px;border:5px solid red" class="w3-round">
    <p>Current channels in server: </p>
  </div>

  <div class="w3-row-padding w3-grayscale" >
    <div class="w3-col m4 w3-margin-bottom" v-for="(c, index) in channels" :key="index" >
      <div class="w3-light-grey" >
        <div class="w3-container">
          <h3>Name: {{c.channel}}</h3>
          <p class="w3-opacity">Members: </p>
          <p v-if="c.members != ''">> {{c.members}}.</p>
          <p v-else>
          > None           </p>
        </div>
      </div>
    </div>
  </div>

  <!-- Payloads -->
  <div class="w3-container" id="designers" style="margin-top:75px">
    <h1 class="w3-xxxlarge w3-text-red"><b>DATA TRANSFER</b></h1>
    <hr style="width:50px;border:5px solid red" class="w3-round">
    <p>Payload bytes per channel: </p>
  </div>

  <div class="w3-row-padding w3-grayscale" >
    <div class="w3-col m4 w3-margin-bottom" v-for="(c, index) in channels" :key="index" >
      <div class="w3-light-grey">
        <div class="w3-container">
          <h3>Name: {{c.channel}}</h3>
          <p class="w3-opacity">Payload bytes: </p>
          <p>{{c.payload}}</p>
        </div>
      </div>
   </div>
  </div>

<!-- End page content -->
</div>

</body>
</html>
</template>

<script>
export default {
 data() {
    return {
      channels: []
    }
  },

  methods: {
    serverupdate: function () {
		this.$http.get('http://localhost:3000/')
    .then(res => this.channels = res.body);
		setTimeout(this.serverupdate, 10000);
	}
  },

  created() {
    this.serverupdate();
  }

}
</script>

<style media= "screen" lang = "css">
  .test{
  background: rgb(255, 255, 255);
  
}
</style>