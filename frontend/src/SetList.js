import React from 'react';
import Table from 'react-bootstrap/Table';

class SetList extends React.Component {

    readData() {
        const self = this;
        fetch("/sets").then(function(response) {
            console.log(response.data);
            self.setState({sets: response.data});
        }).catch(function (error){
            console.log(error);
        });
    }

    getSets() {
        let table = [];

        for (let i=0; i < this.state.sets.length; i++) {
            table.push(
                <tr key={i}>
                    <td>{this.state.sets[i].movement}</td>
                    <td>{this.state.sets[i].volume}</td>
                    <td>{this.state.sets[i].intensity}</td>
                </tr>
            );
        }

        return table;
    }

    constructor(props) {
        super(props);
    }

    componentDidMount() {
        this.readData();
        this.state = {sets: []};
        this.readData = this.readData.bind(this);
    }

    render() {
        return (
            <div>
                <h1 style={{marginBottom: "40px"}}>All Sets</h1>
                <Table>
                    <thead>
                        <tr>
                            <th>
                                Movement
                            </th>
                            <th>
                                Volume
                            </th>
                            <th>
                                Intensity
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        {this.getSets()}
                    </tbody>
                </Table>
            </div>
        )
    }
}

export default SetList;
