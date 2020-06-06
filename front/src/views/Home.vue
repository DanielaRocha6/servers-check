<template>
  <div class="home-view-container">
    <h1>Let's check a domain!</h1>
    <b-form @submit.prevent="handleSubmit">
      <b-form-group id="exampleInputGroup2" >
        <b-form-input
          id="exampleInput2"
          type="text"
          v-model="formData.url"
          required
          placeholder="Enter the URL" />
      </b-form-group>
      <b-button type="submit" variant="primary">Submit</b-button>
    </b-form>
  </div>
</template>

<script>
export default {
  name: 'home',
  data () {
    return {
      formData: {
        url: 'http://'
      }
    }
  },
  methods: {
    alert() {
      this.$swal('Error', 'Please enter a valid website address', 'warning')
    },
    handleSubmit () {
      var { url } = this.formData
      if (this.isUrlValid(url)) {
        if (url.includes('://')) {
          url = url.split('://')[1]
        }
        if (url.includes('www.')) {
          url = url.split('www.')[1]
        }
        const route = '/domains/' + url
        this.$router.push(route)
      } else {
        this.alert()
      }
    },
    isUrlValid(userInput) {
      // eslint-disable-next-line no-useless-escape
      const res = userInput.match(/^((https?|ftp|smtp):\/\/)?(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,4}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)/g)
      return res != null
    }
  }
}
</script>
