import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TablePagination from '@material-ui/core/TablePagination';
import TableRow from '@material-ui/core/TableRow';

function Button(props: any) {
  return (
    <button
      className={props.classname}
      style={{ color: props.textcolor, background: props.backgroundcolor }}
      onClick={props.handle}
    >
      {props.title}
    </button>
  );
}

function Form(props: any) {
  return (
    <label>
      {props.name}
      <input
        className="form"
        type={props.type}
        value={props.value}
        placeholder={props.placeholder}
        onChange={(e) => {
          props.handle(e);
        }}
      />
    </label>
  );
}

const columns = [
  { id: 'topic', label: 'Topic', minWidth: 100 },
  { id: 'payload', label: 'Message', minWidth: 100 },
  { id: 'ts', label: 'Timestamp', minWidth: 100 },
];

const useStyles = makeStyles({
  root: {
    width: '100%',
  },
  container: {
    maxHeight: 440,
  },
});

function StickyHeadTable(props: any) {
  const classes = useStyles();
  const [page, setPage] = React.useState(0);
  const [rowsPerPage, setRowsPerPage] = React.useState(10);

  const handleChangePage = (event: any, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: any) => {
    setRowsPerPage(+event.target.value);
    setPage(0);
  };

  return (
    <Paper className={classes.root}>
      <TableContainer className={classes.container}>
        <Table stickyHeader aria-label="sticky table">
          <TableHead>
            <TableRow>
              {columns.map((column) => (
                <TableCell
                  key={column.id}
                  align={'center'}
                  style={{ minWidth: column.minWidth, backgroundColor: props.backgroundcolor, color: props.textcolor }}
                >
                  {column.label}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {props.rows.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage).map((row: any) => {
              return (
                <TableRow hover role="checkbox" tabIndex={-1} key={row.code}>
                  {columns.map((column) => {
                    const value = row[column.id];
                    return (
                      <TableCell
                        key={column.id}
                        align={'center'}
                        style={{ backgroundColor: props.backgroundcolor, color: props.textcolor }}
                      >
                        {column.id === 'ts' ? new Date(value).toLocaleString() : value}
                      </TableCell>
                    );
                  })}
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        rowsPerPageOptions={props.rowsPerPageOptions}
        component="div"
        count={props.rows.length}
        rowsPerPage={rowsPerPage}
        page={page}
        onChangePage={handleChangePage}
        onChangeRowsPerPage={handleChangeRowsPerPage}
        style={{ backgroundColor: props.backgroundcolor, color: props.textcolor }}
      />
    </Paper>
  );
}

export { Button, Form, StickyHeadTable };
