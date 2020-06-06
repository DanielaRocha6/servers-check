<template>
  <div class="domain-card">
    <b-row class="text-center">
      <b-col cols="7">
        <h5>{{info.Title}}</h5>
        <p class="note">{{url}}</p>
        <img :src="info.Logo" @error="replaceByDefault" :alt="`Logo from ${url}`">

      </b-col>
      <b-col>
        <p><strong>Status: </strong><span :class="getStatusClass()">{{info.Is_down ? ' Down ' : ' Up '}}</span></p>
        <p v-if="info.Ssl_Grade !==''"><strong>Lowest ssl grade: </strong><span :class="getSSLClass('actual')">{{info.Ssl_Grade}}</span></p>
        <p v-if="info.Ssl_Grade ===''"><strong>Lowest ssl grade: </strong><span :class="getSSLClass('actual')">N/A</span></p>

        <p v-if="info.Previous_ssl_grade !==''"><strong>Previous lowest ssl grade: </strong><span :class="getSSLClass('previous')">{{info.Previous_ssl_grade}}</span></p>
        <p v-if="info.Previous_ssl_grade ===''"><strong>Previous lowest ssl grade: </strong><span :class="getSSLClass('previous')">N/A</span></p>
      </b-col>
    </b-row>

  </div>
</template>

<script>
export default {
  name: 'DomainCard',
  props: {
    url: String,
    info: Object
  },
  methods: {
    replaceByDefault(e) {
      e.target.src = `/not-found.png`
    },
    getStatusClass() {
      if (!this.info.Is_down) {
        return 'upStatus'
      } else {
        return 'downStatus'
      }
    },
    getSSLClass(which) {
      const ssl = which === 'actual' ? this.info.Ssl_Grade : this.info.Previous_ssl_grade
      const baseClass = 'ssl-grade'
      if (ssl === 'A+') {
        return baseClass + ' aaplus'
      } else if (ssl === 'A') {
        return baseClass + ' aa'
      } else if (ssl === 'B') {
        return baseClass + ' bb'
      } else if (ssl === 'C') {
        return baseClass + ' cc'
      } else if (ssl === 'D') {
        return baseClass + ' dd'
      } else if (ssl === 'E') {
        return baseClass + ' ee'
      } else if (ssl === 'E') {
        return baseClass + ' ff'
      } else {
        return baseClass + ' na'
      }
    }
  },
  mounted() {
    console.log(this.info)
  }
}
</script>

<style>
.domain-card {
  padding: 1rem;
}
.upStatus {
  background-color: #59B45F;
  border-radius: 15pt;
  padding: 5px;
}
.downStatus {
  background-color: #f76957;
  border-radius: 15pt;
  padding: 5px;
}
</style>
