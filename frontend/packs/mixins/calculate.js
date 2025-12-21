export const calculate = {
  methods: {
    percentProgressBar(current_hp, full_hp) {
      return current_hp / full_hp * 100 + '%'
    }
  }
}
