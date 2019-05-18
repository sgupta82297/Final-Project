class Products extends React.Component {

  constructor(props) {
    super(props);

    //  this.state.products = [];
    this.state = {};
    this.state.filterText = "";
    this.state.products = [
    {
      id: 1,
      category: 'Branch Banking',
      price: '49.99',
      qty: 12,
      name: 'Roger Smith' },
    {
      id: 2,
      category: 'Net Banking',
      price: '9.99',
      qty: 15,
      name: 'Daniel Zarra' },
    {
      id: 3,
      category: 'ATM',
      price: '29.99',
      qty: 14,
      name: 'Chris Hales' },
    {
      id: 4,
      category: 'Zelle Transfer',
      price: '99.99',
      qty: 34,
      name: 'Suzanna Jones' },
    {
      id: 5,
      category: 'Bill Payment',
      price: '399.99',
      qty: 12,
      name: 'Shaun Walter' },
    {
      id: 6,
      category: 'Branch Banking',
      price: '199.99',
      qty: 23,
      name: 'Tom Alter' }];



  }
  handleUserInput(filterText) {
    this.setState({ filterText: filterText });
  }
  handleRowDel(product) {
    var index = this.state.products.indexOf(product);
    this.state.products.splice(index, 1);
    this.setState(this.state.products);
  }

  handleAddEvent(evt) {
    var id = (+new Date() + Math.floor(Math.random() * 999999)).toString(36);
    var product = {
      id: id,
      name: "",
      price: "",
      category: "",
      qty: 0 };

    this.state.products.push(product);
    this.setState(this.state.products);

  }

  handleProductTable(evt) {
    var item = {
      id: evt.target.id,
      name: evt.target.name,
      value: evt.target.value };

    var products = this.state.products.slice();
    var newProducts = products.map(function (product) {

      for (var key in product) {
        if (key == item.name && product.id == item.id) {
          product[key] = item.value;

        }
      }
      return product;
    });
    this.setState({ products: newProducts });
    //  console.log(this.state.products);
  }
  render() {

    return (
      React.createElement("div", null,
      React.createElement(SearchBar, { filterText: this.state.filterText, onUserInput: this.handleUserInput.bind(this) }),
      React.createElement(ProductTable, { onProductTableUpdate: this.handleProductTable.bind(this), onRowAdd: this.handleAddEvent.bind(this), onRowDel: this.handleRowDel.bind(this), products: this.state.products, filterText: this.state.filterText })));



  }}


class SearchBar extends React.Component {
  handleChange() {
    this.props.onUserInput(this.refs.filterTextInput.value);
  }
  render() {
    return (
      React.createElement("div", null,

      React.createElement("input", { type: "text", placeholder: "Search...", value: this.props.filterText, ref: "filterTextInput", onChange: this.handleChange.bind(this) })));




  }}



class ProductTable extends React.Component {

  render() {
    var onProductTableUpdate = this.props.onProductTableUpdate;
    var rowDel = this.props.onRowDel;
    var filterText = this.props.filterText;
    var product = this.props.products.map(function (product) {
      if (product.name.indexOf(filterText) === -1) {
        return;
      }
      return React.createElement(ProductRow, { onProductTableUpdate: onProductTableUpdate, product: product, onDelEvent: rowDel.bind(this), key: product.id });
    });
    return (
      React.createElement("div", null,


      React.createElement("button", { type: "button", onClick: this.props.onRowAdd, className: "btn btn-success pull-right" }, "Add"),
      React.createElement("table", { className: "table table-bordered" },
      React.createElement("thead", null,
      React.createElement("tr", null,
      React.createElement("th", null, "Customer Name"),
      React.createElement("th", null, "Deposit Amount"),
      React.createElement("th", null, "Withdraw Amount"),
      React.createElement("th", null, "Transaction Src"))),



      React.createElement("tbody", null,
      product))));







  }}



class ProductRow extends React.Component {
  onDelEvent() {
    this.props.onDelEvent(this.props.product);

  }
  render() {

    return (
      React.createElement("tr", { className: "eachRow" },
      React.createElement(EditableCell, { onProductTableUpdate: this.props.onProductTableUpdate, cellData: {
          "type": "name",
          value: this.props.product.name,
          id: this.props.product.id } }),

      React.createElement(EditableCell, { onProductTableUpdate: this.props.onProductTableUpdate, cellData: {
          type: "price",
          value: this.props.product.price,
          id: this.props.product.id } }),

      React.createElement(EditableCell, { onProductTableUpdate: this.props.onProductTableUpdate, cellData: {
          type: "qty",
          value: this.props.product.qty,
          id: this.props.product.id } }),

      React.createElement(EditableCell, { onProductTableUpdate: this.props.onProductTableUpdate, cellData: {
          type: "category",
          value: this.props.product.category,
          id: this.props.product.id } }),

      React.createElement("td", { className: "del-cell" },
      React.createElement("input", { type: "button", onClick: this.onDelEvent.bind(this), value: "X", className: "del-btn" }))));




  }}


class EditableCell extends React.Component {

  render() {
    return (
      React.createElement("td", null,
      React.createElement("input", { type: "text", name: this.props.cellData.type, id: this.props.cellData.id, value: this.props.cellData.value, onChange: this.props.onProductTableUpdate })));



  }}


ReactDOM.render(React.createElement(Products, null), document.getElementById('container'));

