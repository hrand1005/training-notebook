import React, { Component } from 'react'
import Table from 'react-bootstrap/Table';

class SetList extends React.Component {

    constructor(){
        super();
        this.state = {
            sets: []
        }
    }

    componentDidMount(){
        fetch('/sets')
        .then(response => response.json())
        .then(data => {
            this.setState({sets: data});
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
