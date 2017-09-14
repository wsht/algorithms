<?php
declare( strict_types = 1 );

function insertion_sort( array $array ) : array
{
	$len = count($array);
	for ($j=1; $j <$len ; $j++) 
	{ 
		$key = $array[$j];

		$i = $j-1;

		while ($i >= 0 && $array[$i] > $key ) 
		{
			$array[$i + 1] = $array[$i];
			$i--;
		}

		$array[$i + 1] = $key;
	}

	return $array;
}

function insertion_sort_desc( array $array ) : array
{
	$len = count( $array );
	for($j = 1; $j < $len ; $j++ )
	{
		$key = $array[$j];
		$i = $j-1;

		while( $i >= 0 && $array[$i] < $key )
		{
			$array[$i + 1] = $array[$i];
			$i--;
		}

		$array[$i + 1] = $key;
	}

	return $array;
}

$array = [123,6,9,66,54,123,421,5,2,1,8,332];

print_r(insertion_sort($array));

print_r(insertion_sort_desc($array));	

