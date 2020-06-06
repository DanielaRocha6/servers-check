<template>
  <div>
    <div >
      <loading :active.sync="isLoading"
               :can-cancel="true"
               :is-full-page="fullPage">
      </loading>
    </div>
    <div v-if="!results">
      <h1>Retrieving SSL information for http://{{url}}</h1>
      <h3>Please wait!</h3>
    </div>
    <div id="domainResult" v-if="results && !resultsError">
      <b-container >
        <b-row>
          <b-col>
            <h4>{{title}}</h4>
            <img :class="getClass()" :src="logo" @error="replaceByDefault" :alt="`Logo from ${url}`">
            <p>http://{{url}}</p>
          </b-col>
          <b-col>
            <b-list-group>
              <b-list-group-item><strong>Status: </strong><span :class="getStatusClass()">{{isDown ? ' Down ' : ' Up '}}</span></b-list-group-item>
              <b-list-group-item v-if="sslGrade!==''"><strong>Lowest ssl grade: </strong><span :class="getSSLClass()">{{sslGrade}}</span></b-list-group-item>
              <b-list-group-item v-if="sslGrade===''"><strong>Lowest ssl grade: </strong><span :class="getSSLClass('actual')">N/A</span></b-list-group-item>
              <b-list-group-item v-if="previousSslGrade!==''"><strong>Previous lowest ssl grade: </strong><span :class="getSSLClass()">{{previousSslGrade}}</span></b-list-group-item>
              <b-list-group-item v-if="previousSslGrade===''"><strong>Previous lowest ssl grade: </strong><span :class="getSSLClass('actual')">N/A</span></b-list-group-item>
            </b-list-group>
          </b-col>
        </b-row>
      </b-container>

      <br>
      <h4>Servers</h4>
      <b-table striped hover :items="actualServers"></b-table>
      <p  class="note">The servers <span><strong>{{serversChanged ? 'changed ' :'had not changed ' }}</strong></span> in the past hour or more</p>
    </div>
    <div v-if="results && resultsError">
      <h2><b>Error found</b></h2>
      <p>The domain {{url}} may not exist or there was an error retrieving the data.</p>
      <h5><a :href="`/domains/${url}`"><b>Please try again!</b></a></h5>
    </div>
  </div>
</template>

<script>
import Loading from 'vue-loading-overlay'
import 'vue-loading-overlay/dist/vue-loading.css'
import axios from 'axios'

export default {
  name: 'domains',
  props: {
    url: {
      type: String,
      default: 'Vue!'
    }
  },
  data () {
    return {
      resultsError: false,
      isLoading: true,
      fullPage: true,
      results: false,
      serversChanged: false,
      sslGrade: '',
      previousSslGrade: '',
      title: '',
      logo: '',
      isDown: '',
      actualServers: [],
      logoErr: false
    }
  },
  methods: {
    getDomain () {
      axios({ method: 'GET', url: 'http://localhost:8090/checkDomain/' + this.url, headers: { 'content-type': 'text/plain' } })
        .then(result => {
          var myResult = result.data
          this.actualServers = myResult.Servers
          if (this.actualServers.length !== 0) {
            this.serversChanged = myResult.Servers_Changed
            this.sslGrade = myResult.Ssl_Grade
            this.previousSslGrade = myResult.Previous_ssl_grade
            this.title = myResult.Title
            this.logo = myResult.Logo
            this.isDown = myResult.Is_down
          } else {
            this.resultsError = true
          }
          this.results = true
          this.isLoading = false
        }).catch(error => {
          this.resultsError = true
          this.results = true
          this.isLoading = false
          console.error(error)
        /* eslint-enable */
        })
    },
    replaceByDefault(e) {
      e.target.src = `/not-found.png`
      this.logoErr = true
    },
    getClass() {
      if (this.logoErr) {
        return 'errClass'
      }
      return ''
    },
    getStatusClass() {
      if (!this.isDown) {
        return 'upStatus'
      } else {
        return 'downStatus'
      }
    },
    getSSLClass(which) {
      const ssl = which === 'actual' ? this.sslGrade : this.previousSslGrade
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
  components: {
    Loading
  },
  mounted () {
    this.getDomain()
  }
}
</script>

<style lang="scss">
  img {
    margin: 1rem;
    width: 5rem;
  }
  .list-group {
    margin-top: 1rem;
    margin-bottom: 1rem;
  }
  .errClass {
    width: 7rem;
  }

</style>
