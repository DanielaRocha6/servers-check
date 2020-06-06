<template>
  <b-navbar toggleable="lg" type="dark" variant="info">
    <b-navbar-brand href="/">Servers check!</b-navbar-brand>

    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

    <b-collapse id="nav-collapse" is-nav>
      <!-- Right aligned nav items -->
      <b-navbar-nav class="ml-auto">
          <b-nav-item>
            <router-link class="text-light" to="/">
              General check
            </router-link>
          </b-nav-item>
        <b-nav-item>
          <router-link class="text-light" to="/domains">
            Domains History
          </router-link>
        </b-nav-item>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>

<script>

export default {
  name: 'Navbar'
}
</script>
<style>

</style>
