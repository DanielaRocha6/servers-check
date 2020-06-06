<template>
  <div id="app">
    <Navbar></Navbar>
    <b-container>
      <router-view/>
    </b-container>
    <Footer></Footer>
  </div>
</template>

<style lang="scss">
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
#nav {
  padding: 30px;
  a {
    font-weight: bold;
    color: #2c3e50;
    &.router-link-exact-active {
      color: #007bff;
    }
  }
}
.container {
  margin-top: 1rem;
}
.aaplus{
  background-color: #59B45F;
}
.aa {
  background-color: #AAD15D;
}
.bb {
  background-color: #E7CD56;
}
.cc {
  background-color: #FCC642;
}
.dd {
  background-color: #FDBA3F;
}
.ee {
  background-color: #F7965A;
}
.ff {
  background-color: #F1624E;
}
.na {
  background-color: lightgray;
}
.ssl-grade {
  border-radius: 15pt;
  padding: 8px;
}
.note {
  font-size: 0.75rem;
}
</style>
<script>
import Navbar from './views/Navbar'
import Footer from './views/Footer'
export default {
  components: { Footer, Navbar }
}
</script>
